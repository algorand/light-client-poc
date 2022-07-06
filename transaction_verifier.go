package main

import (
	"bytes"
	"crypto"
	"fmt"
	"hash"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
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

func getHashedLeaf(txId []byte, stib []byte, hashFunc hash.Hash) []byte {
	buf := make([]byte, 2*hashFunc.Size())
	copy(buf[:], txId[:])
	// TODO: Unsure about the size usage here
	copy(buf[hashFunc.Size():], stib[:])
	leaf := append([]byte{'T', 'L'}, buf...)
	return hashBytes(hashFunc, leaf)
}

func getVectorCommitmentPositions(index uint64, depth uint64) []Position {
	directions := make([]Position, depth)
	for i := len(directions) - 1; i >= 0; i-- {
		directions[i] = Position(index & 1)
		index >>= 1
	}

	return directions
}

func verifyTransaction(transactionCommitment, txId []byte, transactionProof models.ProofResponse) (bool, error) {
	hashFunc, err := unmarshalHashFunc(transactionProof.Hashtype)
	if err != nil {
		return false, err
	}

	nodeSize := uint64(hashFunc.Size())
	currentNodeHash := getHashedLeaf(txId, transactionProof.Stibhash, hashFunc)

	positions := getVectorCommitmentPositions(transactionProof.Idx, transactionProof.Treedepth)
	for i := uint64(0); i < transactionProof.Treedepth; i++ {
		siblingIndex := i * nodeSize
		siblingHash := transactionProof.Proof[siblingIndex : siblingIndex+nodeSize]

		nextNode := []byte{'M', 'A'}
		switch positions[i] {
		case left:
			nextNode = append(append(nextNode, currentNodeHash...), siblingHash...)
		case right:
			nextNode = append(append(nextNode, siblingHash...), currentNodeHash...)
		default:
			return false, fmt.Errorf("bad direction")
		}

		currentNodeHash = hashBytes(hashFunc, nextNode)
	}

	return bytes.Equal(currentNodeHash, transactionCommitment), nil
}
