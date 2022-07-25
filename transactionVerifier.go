package main

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"hash"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/stateproofs/stateprooftypes"
	"github.com/algorand/go-algorand-sdk/types"
)

var (
	ErrUnsupportedHashFunction      = errors.New("proof hash function is unsupported")
	ErrProofLengthTreeDepthMismatch = errors.New("proof length and tree depth do not match")
	ErrRootMismatch                 = errors.New("root mismatch")
	ErrInvalidTreeDepth             = errors.New("invalid tree depth")
	ErrIndexDepthMismatch           = errors.New("node index is not smaller than 2^depth")
	ErrInvalidPosition              = errors.New("invalid position for node")
)

type Position int

const (
	left Position = iota
	right
)

type TransactionVerifier struct {
	genesisHash types.Digest
}

func (t *TransactionVerifier) getTransactionLeaf(txId types.Digest, stib types.Digest, hashFunc hash.Hash) stateprooftypes.GenericDigest {
	buf := make([]byte, 2*types.DigestSize)
	copy(buf[:], txId[:])
	copy(buf[types.DigestSize:], stib[:])
	leaf := append([]byte(stateprooftypes.TxnMerkleLeaf), buf...)
	return stateprooftypes.HashBytes(hashFunc, leaf)
}

func (t *TransactionVerifier) getLightBlockHeaderLeaf(roundNumber types.Round,
	transactionCommitment stateprooftypes.GenericDigest, seed stateprooftypes.Seed, hashFunc hash.Hash) stateprooftypes.GenericDigest {
	lightBlockheader := stateprooftypes.LightBlockHeader{
		RoundNumber:         roundNumber,
		GenesisHash:         t.genesisHash,
		Sha256TxnCommitment: transactionCommitment,
		Seed:                seed,
	}

	lightBlockheader.ToBeHashed()
	return stateprooftypes.HashBytes(hashFunc, lightBlockheader.ToBeHashed())
}

// getVectorCommitmentPositions takes the index provided with the proof, which is a bitmap of positions,
// and translates it to an array of positions, starting from the LSB - this is the key difference
// between vector commitments and merkle trees. The resulting array should contain a number of positions equal to the
// depth of the tree - the length of the path between the root and the leaves.
// Since each bit of the index is mapped to an element in the resulting positions array,
// the index must contain a number of bits amounting to the depth parameter, which means it must be smaller than 2 ^ depth.
func getVectorCommitmentPositions(index uint64, depth uint64) ([]Position, error) {
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

func computeVectorCommitmentMerkleRoot(leaf stateprooftypes.GenericDigest, leafIndex uint64, proof []byte, treeDepth uint64, hashFunc hash.Hash) (stateprooftypes.GenericDigest, error) {
	if len(proof) == 0 && treeDepth == 0 {
		return leaf, nil
	}

	nodeSize := uint64(hashFunc.Size())
	if treeDepth*nodeSize != uint64(len(proof)) {
		return stateprooftypes.GenericDigest{}, ErrProofLengthTreeDepthMismatch
	}

	positions, err := getVectorCommitmentPositions(leafIndex, treeDepth)
	if err != nil {
		return stateprooftypes.GenericDigest{}, err
	}

	currentNodeHash := leaf
	for currentLevel := uint64(0); currentLevel < treeDepth; currentLevel++ {
		siblingIndex := currentLevel * nodeSize
		siblingHash := proof[siblingIndex : siblingIndex+nodeSize]

		nextNode := []byte(stateprooftypes.MerkleArrayNode)
		switch positions[currentLevel] {
		case left:
			nextNode = append(append(nextNode, currentNodeHash...), siblingHash...)
		case right:
			nextNode = append(append(nextNode, siblingHash...), currentNodeHash...)
		default:
			return []byte{}, ErrInvalidPosition
		}

		currentNodeHash = stateprooftypes.HashBytes(hashFunc, nextNode)
	}

	return currentNodeHash, nil
}

// VerifyTransaction receives a sha256 hashed transaction, a proof to compute the transaction's commitment, a proof
// to compute the light block header's commitment, and an expected commitment to compare to.
// Verification works as follows:
//	1. Compute the transaction's vector commitment using the provided transaction proof.
//	2. Build a candidate light block header using the computed vector commitment.
// 	3. Compute the candidate light block header's vector commitment using the provided light block header proof.
// 	4. Verify that the computed candidate vector commitment matches the expected vector commitment.
func (t *TransactionVerifier) VerifyTransaction(transactionId types.Digest, transactionProofResponse models.ProofResponse, lightBlockHeaderProofResponse models.LightBlockHeaderProof, confirmedRound types.Round, seed stateprooftypes.Seed, blockIntervalCommitment types.Digest) error {
	// verifying attested vector commitments is currently exclusively supported with sha256 hashing, both for transactions
	// and light block headers.
	if transactionProofResponse.Hashtype != "sha256" {
		return ErrUnsupportedHashFunction
	}

	var stibHashDigest types.Digest
	copy(stibHashDigest[:], transactionProofResponse.Stibhash[:])

	transactionHashFunc := sha256.New()
	transactionLeaf := t.getTransactionLeaf(transactionId, stibHashDigest, transactionHashFunc)
	transactionProofRoot, err := computeVectorCommitmentMerkleRoot(transactionLeaf, transactionProofResponse.Idx,
		transactionProofResponse.Proof, transactionProofResponse.Treedepth, transactionHashFunc)

	if err != nil {
		return err
	}

	lightBlockHeaderHashFunc := sha256.New()
	// We build the candidate light block header using the computed transactionProofRoot, hash and verify it.
	candidateLightBlockHeaderLeaf := t.getLightBlockHeaderLeaf(confirmedRound, transactionProofRoot, seed, lightBlockHeaderHashFunc)
	lightBlockHeaderProofRoot, err := computeVectorCommitmentMerkleRoot(candidateLightBlockHeaderLeaf, lightBlockHeaderProofResponse.Index, lightBlockHeaderProofResponse.Proof,
		lightBlockHeaderProofResponse.Treedepth, lightBlockHeaderHashFunc)

	if err != nil {
		return err
	}

	if bytes.Equal(lightBlockHeaderProofRoot, blockIntervalCommitment[:]) != true {
		return ErrRootMismatch
	}
	return nil
}
