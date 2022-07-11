package main

import (
	"bytes"
	"crypto"
	"fmt"
	"hash"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
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

func getTransactionLeaf(txId []byte, stib []byte, hashFunc hash.Hash) []byte {
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

func VerifyTransaction(transactionId []byte, transactionProofResponse models.ProofResponse,
	lightBlockHeaderProofResponse LightBlockHeaderProofResponse, blockIntervalCommitment []byte, genesisHash []byte, roundNumber uint64) (bool, error) {
	hashFunc, err := unmarshalHashFunc(transactionProofResponse.Hashtype)
	if err != nil {
		return false, err
	}

	transactionLeaf := getTransactionLeaf(transactionId, transactionProofResponse.Stibhash, hashFunc)
	transactionProofRoot, err := climbProof(transactionLeaf, transactionProofResponse.Idx,
		transactionProofResponse.Proof, transactionProofResponse.Treedepth, hashFunc)

	if err != nil {
		return false, err
	}

	lightBlockHeaderLeaf := getLightBlockHeaderLeaf(genesisHash, roundNumber, transactionProofRoot, hashFunc)
	lightBlockHeaderProofRoot, err := climbProof(lightBlockHeaderLeaf, lightBlockHeaderProofResponse.Index, lightBlockHeaderProofResponse.Proof,
		lightBlockHeaderProofResponse.Treedepth, hashFunc)

	if err != nil {
		return false, err
	}

	return bytes.Equal(lightBlockHeaderProofRoot, blockIntervalCommitment), nil
}
