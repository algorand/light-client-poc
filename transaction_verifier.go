package main

import (
	"bytes"
	"crypto"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
	"hash"
)

type Position int

const (
	left Position = iota
	right
)

func hashBytes(hash hash.Hash, m []byte) []byte {
	hash.Reset()
	hash.Write(m)
	outhash := hash.Sum(nil)
	return outhash
}

func unmarshalHashFunc(hashStr string) (hash.Hash, error) {
	switch hashStr {
	case "sha256":
		return crypto.SHA256.New(), nil
	default:
		return nil, fmt.Errorf("unsupported hash function detected")
	}
}

func getHashedTransactionLeaf(txId []byte, stib []byte, hashFunc hash.Hash) []byte {
	buf := make([]byte, 2*hashFunc.Size())
	copy(buf[:], txId[:])
	// TODO: Unsure about the size usage here
	copy(buf[hashFunc.Size():], stib[:])
	leaf := append([]byte{'T', 'L'}, buf...)
	return hashBytes(hashFunc, leaf)
}

func getLightBlockHeaderLeaf(genesisHash []byte, roundNumber uint64, transactionCommitment []byte, hashFunc hash.Hash) []byte {
	lightBlockheader := types.LightBlockHeader{
		RoundNumber:         roundNumber,
		GenesisHash:         genesisHash,
		Sha256TxnCommitment: transactionCommitment,
	}

	lightBlockheader.ToBeHashed()
	return hashBytes(hashFunc, lightBlockheader.ToBeHashed())
}

func getVectorCommitmentPositions(index uint64, depth uint64) []Position {
	directions := make([]Position, depth)
	for i := len(directions) - 1; i >= 0; i-- {
		directions[i] = Position(index & 1)
		index >>= 1
	}

	return directions
}

func climbProof(leaf []byte, leafIndex uint64, proof []byte, treeDepth uint64, hashFunc hash.Hash) ([]byte, error) {

	nodeSize := uint64(hashFunc.Size())
	//currentNodeHash := getHashedTransactionLeaf(txId, transactionProof.Stibhash, hashFunc)
	currentNodeHash := leaf

	positions := getVectorCommitmentPositions(leafIndex, treeDepth)
	for i := uint64(0); i < treeDepth; i++ {
		siblingIndex := i * nodeSize
		siblingHash := proof[siblingIndex : siblingIndex+nodeSize]

		nextNode := []byte{'M', 'A'}
		switch positions[i] {
		case left:
			nextNode = append(append(nextNode, currentNodeHash...), siblingHash...)
		case right:
			nextNode = append(append(nextNode, siblingHash...), currentNodeHash...)
		default:
			return []byte{}, fmt.Errorf("bad direction")
		}

		currentNodeHash = hashBytes(hashFunc, nextNode)
	}

	return currentNodeHash, nil
}

func getTransactionProofRoot(transactionId []byte, transactionProof models.ProofResponse) ([]byte, error) {
	hashFunc, err := unmarshalHashFunc(transactionProof.Hashtype)
	if err != nil {
		return []byte{}, err
	}

	hashedLeaf := getHashedTransactionLeaf(transactionId, transactionProof.Stibhash, hashFunc)
	return climbProof(hashedLeaf, transactionProof.Idx, transactionProof.Proof, transactionProof.Treedepth, hashFunc)
}

func VerifyTransaction(transactionId []byte, transactionProof models.ProofResponse,
	lightBlockHeaderProof LightBlockHeaderProof, blockIntervalCommitment []byte, genesisHash []byte, roundNumber uint64) (bool, error) {
	transactionProofRoot, err := getTransactionProofRoot(transactionId, transactionProof)
	if err != nil {
		return false, err
	}

	sha256HashFunc := crypto.SHA256.New()
	lightBlockHeaderLeaf := getLightBlockHeaderLeaf(genesisHash, roundNumber, transactionProofRoot, sha256HashFunc)
	lightBlockHeaderProofRoot, err := climbProof(lightBlockHeaderLeaf, lightBlockHeaderProof.Index, lightBlockHeaderProof.Proof,
		lightBlockHeaderProof.Treedepth, sha256HashFunc)

	if err != nil {
		return false, err
	}

	return bytes.Equal(lightBlockHeaderProofRoot, blockIntervalCommitment), nil
}
