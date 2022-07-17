package main

import (
	"bytes"
	"fmt"
	"hash"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/stateproofs/datatypes"

	"light-client-poc/encoded_assets"
	"light-client-poc/utilities"
)

type Position int

const (
	left Position = iota
	right
)

type TransactionVerifier struct {
	genesisHash datatypes.GenericDigest
}

func (t *TransactionVerifier) getTransactionLeaf(txId []byte, stib []byte, hashFunc hash.Hash) []byte {
	buf := make([]byte, 2*hashFunc.Size())
	copy(buf[:], txId[:])
	// TODO: Unsure about the size usage here
	copy(buf[hashFunc.Size():], stib[:])
	leaf := append([]byte{'T', 'L'}, buf...)
	return utilities.HashBytes(hashFunc, leaf)
}

func (t *TransactionVerifier) getLightBlockHeaderLeaf(roundNumber uint64, transactionCommitment []byte, hashFunc hash.Hash) []byte {
	lightBlockheader := datatypes.LightBlockHeader{
		RoundNumber:         roundNumber,
		GenesisHash:         t.genesisHash,
		Sha256TxnCommitment: transactionCommitment,
	}

	lightBlockheader.ToBeHashed()
	return utilities.HashBytes(hashFunc, lightBlockheader.ToBeHashed())
}

func (t *TransactionVerifier) getVectorCommitmentPositions(index uint64, depth uint64) []Position {
	directions := make([]Position, depth)
	for i := len(directions) - 1; i >= 0; i-- {
		directions[i] = Position(index & 1)
		index >>= 1
	}

	return directions
}

func (t *TransactionVerifier) climbProof(leaf []byte, leafIndex uint64, proof []byte, treeDepth uint64, hashFunc hash.Hash) ([]byte, error) {
	nodeSize := uint64(hashFunc.Size())
	currentNodeHash := leaf

	// TODO: Verify proof according to node size

	positions := t.getVectorCommitmentPositions(leafIndex, treeDepth)
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

		currentNodeHash = utilities.HashBytes(hashFunc, nextNode)
	}

	return currentNodeHash, nil
}

// TODO: Change to ptr
func (t *TransactionVerifier) VerifyTransaction(transactionId []byte, transactionProofResponse models.ProofResponse,
	lightBlockHeaderProofResponse encoded_assets.LightBlockHeaderProofResponse, blockIntervalCommitment []byte, roundNumber uint64) (bool, error) {
	hashFunc, err := utilities.UnmarshalHashFunc(transactionProofResponse.Hashtype)
	if err != nil {
		return false, err
	}

	transactionLeaf := t.getTransactionLeaf(transactionId, transactionProofResponse.Stibhash, hashFunc)
	transactionProofRoot, err := t.climbProof(transactionLeaf, transactionProofResponse.Idx,
		transactionProofResponse.Proof, transactionProofResponse.Treedepth, hashFunc)

	// TODO: Add verification of transactionProofRoot with inputted SHA256TxnRoot? Is it necessary, considering that
	// TODO: SHA256TxnRoot is not trusted itself?
	if err != nil {
		return false, err
	}

	lightBlockHeaderLeaf := t.getLightBlockHeaderLeaf(roundNumber, transactionProofRoot, hashFunc)
	lightBlockHeaderProofRoot, err := t.climbProof(lightBlockHeaderLeaf, lightBlockHeaderProofResponse.Index, lightBlockHeaderProofResponse.Proof,
		lightBlockHeaderProofResponse.Treedepth, hashFunc)

	if err != nil {
		return false, err
	}

	return bytes.Equal(lightBlockHeaderProofRoot, blockIntervalCommitment), nil
}
