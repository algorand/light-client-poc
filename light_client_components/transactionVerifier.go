package light_client_components

import (
	"bytes"
	"crypto/sha256"
	"errors"
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

type NodePosition int

const (
	leftSibling NodePosition = iota
	rightSibling
)

// computeTransactionLeaf receives the transaction ID and the signed transaction in block's hash, and computes
// the leaf of the vector commitment associated with the transaction.
// Parameters:
// txId is the Sha256 hash of the transaction.
// stibHash is the Sha256 of the transaction as it's saved in the block.
func computeTransactionLeaf(txId types.Digest, stibHash types.Digest) types.Digest {
	buf := make([]byte, 2*types.DigestSize)
	copy(buf[:], txId[:])
	copy(buf[types.DigestSize:], stibHash[:])
	leaf := append([]byte(stateprooftypes.TxnMerkleLeaf), buf...)
	// The leaf returned is of the form: Sha256("TL", Sha256(transaction), Sha256(transaction in block))
	return sha256.Sum256(leaf)
}

// computeLightBlockHeaderLeaf receives the parameters comprising a light block header, and computes the leaf
// of the vector commitment associated with the transaction.
// Parameters:
// roundNumber is the round of the block to which the light block header belongs.
// transactionCommitment is the sha256 vector commitment for the transactions in the block to which the light block header belongs.
// genesisHash is the hash of the genesis block.
// seed is the block's sortition seed.
func computeLightBlockHeaderLeaf(roundNumber types.Round,
	transactionCommitment types.Digest, genesisHash types.Digest, seed stateprooftypes.Seed) types.Digest {
	lightBlockheader := stateprooftypes.LightBlockHeader{
		RoundNumber:         roundNumber,
		GenesisHash:         genesisHash,
		Sha256TxnCommitment: transactionCommitment,
		Seed:                seed,
	}

	// The leaf returned is of the form Sha256(lightBlockHeader)
	return sha256.Sum256(lightBlockheader.ToBeHashed())
}

// getVectorCommitmentPositions takes the index provided with the proof, which is a bitmap of positions,
// and translates it to an array of node positions, starting from the LSB - this is the key difference
// between vector commitments and merkle trees. Each position at index i of the result
// corresponds to our expected node position relative to its sibling when computing the root,
// at height i relative to the leaf.
// For example, positions[0] is the position of the leaf relative to its sibling.
// Parameters:
// index is the node's index in the vector commitment tree.
// depth is the node's depth in the vector commitment tree.
func getVectorCommitmentPositions(index uint64, depth uint64) ([]NodePosition, error) {
	// A depth of 0 is only valid when the proof's length is also 0. Since the calling functions checks for the situation
	// where both values are 0, depth of 0 must be invalid.
	if depth == 0 {
		return []NodePosition{}, ErrInvalidTreeDepth
	}

	// Since each bit of the index is mapped to an element in the resulting positions array,
	// the index must contain a number of bits amounting to the depth parameter, which means it must be smaller than 2 ^ depth.
	if index >= 1<<depth {
		return []NodePosition{}, ErrIndexDepthMismatch
	}

	// The resulting array should contain a number of positions equal to the depth of the tree -
	// the length of the path between the root and the leaves - as that is the amounts of nodes traversed when calculating
	// the vector commitment root.
	directions := make([]NodePosition, depth)

	// We iterate on the resulting array starting from the end, to allow us to extract LSBs, yet have the eventual result
	// be equivalent to extracting MSBs.
	for i := len(directions) - 1; i >= 0; i-- {
		// We take index's current LSB, translate it to a node position and place it in index i.
		directions[i] = NodePosition(index & 1)
		// We shift index to the right, to prepare for the extracting of the next LSB.
		index >>= 1
	}

	return directions, nil
}

func computeVectorCommitmentMerkleRoot(leaf types.Digest, leafIndex uint64, proof []byte, treeDepth uint64) (types.Digest, error) {
	if len(proof) == 0 && treeDepth == 0 {
		return leaf, nil
	}

	nodeSize := uint64(sha256.New().Size())
	if treeDepth*nodeSize != uint64(len(proof)) {
		return types.Digest{}, ErrProofLengthTreeDepthMismatch
	}

	positions, err := getVectorCommitmentPositions(leafIndex, treeDepth)
	if err != nil {
		return types.Digest{}, err
	}

	currentNodeHash := leaf
	for currentLevel := uint64(0); currentLevel < treeDepth; currentLevel++ {
		siblingIndex := currentLevel * nodeSize
		siblingHash := proof[siblingIndex : siblingIndex+nodeSize]

		nextNode := []byte(stateprooftypes.MerkleArrayNode)
		switch positions[currentLevel] {
		case leftSibling:
			nextNode = append(append(nextNode, currentNodeHash[:]...), siblingHash...)
		case rightSibling:
			nextNode = append(append(nextNode, siblingHash...), currentNodeHash[:]...)
		default:
			return types.Digest{}, ErrInvalidPosition
		}

		currentNodeHash = sha256.Sum256(nextNode)
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
func VerifyTransaction(transactionHash types.Digest, transactionProofResponse models.ProofResponse,
	lightBlockHeaderProofResponse models.LightBlockHeaderProof, confirmedRound types.Round, genesisHash types.Digest, seed stateprooftypes.Seed, blockIntervalCommitment types.Digest) error {
	// verifying attested vector commitments is currently exclusively supported with sha256 hashing, both for transactions
	// and light block headers.
	if transactionProofResponse.Hashtype != "sha256" {
		return ErrUnsupportedHashFunction
	}

	var stibHashDigest types.Digest
	copy(stibHashDigest[:], transactionProofResponse.Stibhash[:])

	transactionLeaf := computeTransactionLeaf(transactionHash, stibHashDigest)
	transactionProofRoot, err := computeVectorCommitmentMerkleRoot(transactionLeaf, transactionProofResponse.Idx,
		transactionProofResponse.Proof, transactionProofResponse.Treedepth)

	if err != nil {
		return err
	}

	// We build the candidate light block header using the computed transactionProofRoot, hash and verify it.
	candidateLightBlockHeaderLeaf := computeLightBlockHeaderLeaf(confirmedRound, transactionProofRoot, genesisHash, seed)
	lightBlockHeaderProofRoot, err := computeVectorCommitmentMerkleRoot(candidateLightBlockHeaderLeaf, lightBlockHeaderProofResponse.Index, lightBlockHeaderProofResponse.Proof,
		lightBlockHeaderProofResponse.Treedepth)

	if err != nil {
		return err
	}

	if bytes.Equal(lightBlockHeaderProofRoot[:], blockIntervalCommitment[:]) != true {
		return ErrRootMismatch
	}
	return nil
}
