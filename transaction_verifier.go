package main

import (
	"bytes"
	"errors"
	"hash"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/stateproofs/datatypes"
	"github.com/algorand/go-algorand-sdk/types"
)

var (
	ErrProofLengthTreeDepthMismatch = errors.New("proof length and tree depth do not match")
	ErrRootMismatch                 = errors.New("root mismatch")
	ErrInvalidTreeDepth             = errors.New("invalid tree depth")
	ErrIndexDepthMismatch           = errors.New("node index is not smaller than 2^depth")
	ErrInvalidPosition              = errors.New("invalid position for node")
)

type Position int

const (
	right Position = iota
	left
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

func (t *TransactionVerifier) getVectorCommitmentPositions(index uint64, depth uint64) ([]Position, error) {
	if depth == 0 {
		return []Position{}, ErrInvalidTreeDepth
	}

	if index >= 1<<depth {
		return []Position{}, ErrIndexDepthMismatch
	}

	directions := make([]Position, depth)
	for i := len(directions) - 1; i >= 0; i-- {
		directions[i] = Position(index & 1)
		index >>= 1
	}

	return directions, nil
}

func (t *TransactionVerifier) computeMerkleRoot(leaf datatypes.GenericDigest, leafIndex uint64, proof []byte, treeDepth uint64, hashFunc hash.Hash) (datatypes.GenericDigest, error) {
	if len(proof) == 0 && treeDepth == 0 {
		return leaf, nil
	}

	nodeSize := uint64(hashFunc.Size())
	if treeDepth*nodeSize != uint64(len(proof)) {
		return datatypes.GenericDigest{}, ErrProofLengthTreeDepthMismatch
	}
	
	positions, err := t.getVectorCommitmentPositions(leafIndex, treeDepth)
	if err != nil {
		return datatypes.GenericDigest{}, err
	}

	currentNodeHash := leaf
	for i := uint64(0); i < treeDepth; i++ {
		siblingIndex := i * nodeSize
		siblingHash := proof[siblingIndex : siblingIndex+nodeSize]

		nextNode := []byte(datatypes.MerkleArrayNode)
		switch positions[i] {
		case right:
			nextNode = append(append(nextNode, currentNodeHash...), siblingHash...)
		case left:
			nextNode = append(append(nextNode, siblingHash...), currentNodeHash...)
		default:
			return []byte{}, ErrInvalidPosition
		}

		currentNodeHash = datatypes.HashBytes(hashFunc, nextNode)
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

	if bytes.Equal(lightBlockHeaderProofRoot, blockIntervalCommitment[:]) != true {
		return ErrRootMismatch
	}
	return nil
}
