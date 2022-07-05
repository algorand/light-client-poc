package main

import (
	"bytes"
	"crypto"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"hash"
)

type Direction int

const (
	left Direction = iota
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

func getPreimage(txId []byte, stib []byte) []byte {
	buf := make([]byte, 64)
	copy(buf[:], txId[:])
	copy(buf[32:], stib[:])
	return append([]byte{'T', 'L'}, buf...)
}

func verifyTransaction(commitment, txId []byte, transactionProof models.ProofResponse) bool {
	nodeSize := uint64(32)

	hashFunc, err := unmarshalHashFunc(transactionProof.Hashtype)
	if err != nil {
		return false
	}

	preImage := getPreimage(txId, transactionProof.Stibhash)
	currentNodeHash := hashBytes(hashFunc, preImage)

	//idxDirection := bits.Reverse64(transactionProof.Idx)
	directions := getPathDirections(transactionProof.Idx, transactionProof.Treedepth)
	for i := uint64(0); i < transactionProof.Treedepth; i++ {
		currentNodeIdx := i * nodeSize
		siblingHash := transactionProof.Proof[currentNodeIdx : currentNodeIdx+nodeSize]

		nextNode := []byte{'M', 'A'}
		if directions[i] == left {
			nextNode = append(append(nextNode, currentNodeHash...), siblingHash...)
		} else {
			nextNode = append(append(nextNode, siblingHash...), currentNodeHash...)
		}

		currentNodeHash = hashBytes(hashFunc, nextNode)
	}

	return bytes.Equal(currentNodeHash, commitment)
}

func getPathDirections(index uint64, depth uint64) []Direction {
	directions := make([]Direction, depth)
	for i, _ := range directions {
		directions[i] = Direction(index & 1)
		index = index >> 1
	}

	return directions
}
