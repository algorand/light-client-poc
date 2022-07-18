package main

import (
	"bytes"
	"fmt"
	"hash"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/stateproofs/datatypes"
	"github.com/algorand/go-algorand-sdk/types"

	"light-client-poc/utilities"
)

type Position int

const (
	left Position = iota
	right
)

type TransactionVerifier struct {
	genesisHash types.Digest
}

func (t *TransactionVerifier) getTransactionLeaf(txId types.Digest, stib types.Digest, hashFunc hash.Hash) datatypes.GenericDigest {
	buf := make([]byte, 2*types.DigestSize)
	copy(buf[:], txId[:])
	copy(buf[types.DigestSize:], stib[:])
	leaf := append([]byte(datatypes.TxnMerkleLeaf), buf...)
	return datatypes.HashBytes(hashFunc, leaf)
}

func (t *TransactionVerifier) getLightBlockHeaderLeaf(roundNumber types.Round, transactionCommitment datatypes.GenericDigest, hashFunc hash.Hash) datatypes.GenericDigest {
	lightBlockheader := datatypes.LightBlockHeader{
		RoundNumber:         roundNumber,
		GenesisHash:         t.genesisHash,
		Sha256TxnCommitment: transactionCommitment,
	}

	lightBlockheader.ToBeHashed()
	return datatypes.HashBytes(hashFunc, lightBlockheader.ToBeHashed())
}

func (t *TransactionVerifier) getVectorCommitmentPositions(index uint64, depth uint64) []Position {
	directions := make([]Position, depth)
	for i := len(directions) - 1; i >= 0; i-- {
		directions[i] = Position(index & 1)
		index >>= 1
	}

	return directions
}

func (t *TransactionVerifier) computeMerkleRoot(leaf datatypes.GenericDigest, leafIndex uint64, proof []byte, treeDepth uint64, hashFunc hash.Hash) (datatypes.GenericDigest, error) {
	nodeSize := uint64(hashFunc.Size())
	currentNodeHash := leaf

	// TODO: Verify proof according to node size

	positions := t.getVectorCommitmentPositions(leafIndex, treeDepth)
	for i := uint64(0); i < treeDepth; i++ {
		siblingIndex := i * nodeSize
		siblingHash := proof[siblingIndex : siblingIndex+nodeSize]

		nextNode := []byte(datatypes.MerkleArrayNode)
		switch positions[i] {
		case left:
			nextNode = append(append(nextNode, currentNodeHash...), siblingHash...)
		case right:
			nextNode = append(append(nextNode, siblingHash...), currentNodeHash...)
		default:
			// TODO: Create error type
			return []byte{}, fmt.Errorf("bad direction")
		}

		currentNodeHash = utilities.HashBytes(hashFunc, nextNode)
	}

	return currentNodeHash, nil
}

func (t *TransactionVerifier) VerifyTransaction(transactionId types.Digest, transactionProofResponse models.ProofResponse,
	lightBlockHeaderProofResponse models.LightBlockHeaderProof, blockIntervalCommitment types.Digest, roundNumber types.Round) error {
	hashFunc, err := datatypes.UnmarshalHashFunc(transactionProofResponse.Hashtype)
	if err != nil {
		return err
	}

	var stibHashDigest types.Digest
	copy(stibHashDigest[:], transactionProofResponse.Stibhash[:])
	transactionLeaf := t.getTransactionLeaf(transactionId, stibHashDigest, hashFunc)
	transactionProofRoot, err := t.computeMerkleRoot(transactionLeaf, transactionProofResponse.Idx,
		transactionProofResponse.Proof, transactionProofResponse.Treedepth, hashFunc)

	// TODO: Add verification of transactionProofRoot with inputted SHA256TxnRoot? Is it necessary, considering that
	// TODO: SHA256TxnRoot is not trusted itself?
	if err != nil {
		return err
	}

	lightBlockHeaderLeaf := t.getLightBlockHeaderLeaf(roundNumber, transactionProofRoot, hashFunc)
	lightBlockHeaderProofRoot, err := t.computeMerkleRoot(lightBlockHeaderLeaf, lightBlockHeaderProofResponse.Index, lightBlockHeaderProofResponse.Proof,
		lightBlockHeaderProofResponse.Treedepth, hashFunc)

	if err != nil {
		return err
	}

	// TODO: Add error type
	if bytes.Equal(lightBlockHeaderProofRoot, blockIntervalCommitment[:]) != true {
		return fmt.Errorf("calculated root and trusted commitment are different")
	}
	return nil
}
