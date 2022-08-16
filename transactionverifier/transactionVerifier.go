package transactionverifier

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/stateproofs/transactionverificationtypes"
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
	leftChild NodePosition = iota
	rightChild
)

// computeTransactionLeaf receives the transaction ID and the signed transaction in block's hash, and computes
// the leaf of the vector commitment associated with the transaction.
// Parameters:
// txId - the Sha256 hash of the msgpacked transaction.
// stibHash - the Sha256 of the msgpacked transaction as it's saved in the block.
func computeTransactionLeaf(txId types.Digest, stibHash types.Digest) types.Digest {
	leafDomainSeparator := []byte(transactionverificationtypes.TxnMerkleLeaf)

	var leafData []byte
	leafData = append(leafData, leafDomainSeparator...)
	leafData = append(leafData, txId[:]...)
	leafData = append(leafData, stibHash[:]...)

	// The leaf returned is of the form: Sha256("TL" || Sha256(transaction) || Sha256(transaction in block))
	return sha256.Sum256(leafData)
}

// computeLightBlockHeaderLeaf receives the parameters comprising a light block header, and computes the leaf
// of the vector commitment associated with the transaction.
// Parameters:
// roundNumber - the round of the block to which the light block header belongs.
// transactionCommitment - the sha256 vector commitment root for the transactions in the block to which the light block header belongs.
// genesisHash - the hash of the genesis block.
// seed - the sortition seed of the block associated with the light block header.
func computeLightBlockHeaderLeaf(roundNumber types.Round,
	transactionCommitment types.Digest, genesisHash types.Digest, seed transactionverificationtypes.Seed) types.Digest {
	lightBlockheader := transactionverificationtypes.LightBlockHeader{
		RoundNumber:         roundNumber,
		GenesisHash:         genesisHash,
		Sha256TxnCommitment: transactionCommitment,
		Seed:                seed,
	}

	// The leaf returned is of the form Sha256(lightBlockHeader)
	return sha256.Sum256(lightBlockheader.ToBeHashed())
}

// getVectorCommitmentPositions maps a depth and a vector commitment index to the "positions" of the nodes
// on the leaf-to-root path, with 0/1 denoting left/right (respectively). It does so by expressing the index in binary
// using exactly depth bits, starting from the most-significant bit.
// For example, leafDepth=4 and index=5 maps to the array [0,1,0,1] -- indicating that the leaf is a left child,
// the leaf's parent is a right child, etc.
// Parameters:
// index - the leaf's index in the vector commitment tree.
// depth - the length of the path from the leaf to the root.
func getVectorCommitmentPositions(index uint64, depth uint64) ([]NodePosition, error) {
	// A depth of 0 is only valid when the proof's length is also 0. Since the calling function checks for the situation
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
		// We shift the index to the right, to prepare for the extracting of the next LSB.
		index >>= 1
	}

	return directions, nil
}

// computeVectorCommitmentRoot takes a vector commitment leaf, its index, a proof, and a tree depth. it calculates
// the vector commitment root using the provided data. This is done by computing internal nodes using the proof,
// starting from the leaf, until we reach the root. This function uses sha256 - it cannot be used correctly with leaves
// and proofs created using a different hash function.
// Parameters:
// leaf - the node we start computing the vector commitment root from.
// leafIndex - the leaf's index.
// proof - the proof to use in computing the vector commitment root. It holds hashed sibling nodes for each internal node
// calculated.
// treeDepth - the length of the path from the leaf to the root.
func computeVectorCommitmentRoot(leaf types.Digest, leafIndex uint64, proof []byte, treeDepth uint64) (types.Digest, error) {
	// An empty proof is only possible when the leaf received is already the root, which means that the treeDepth
	// must be 0. In this case, the result is the leaf itself.
	if len(proof) == 0 && treeDepth == 0 {
		return leaf, nil
	}

	nodeHashSize := uint64(sha256.New().Size())
	// The proof must hold exactly treeDepth node hashes to allow us to compute enough nodes to reach the root.
	if treeDepth*nodeHashSize != uint64(len(proof)) {
		return types.Digest{}, ErrProofLengthTreeDepthMismatch
	}

	// See comments on getVectorCommitmentPositions for more details on the contents of the positions variable.
	positions, err := getVectorCommitmentPositions(leafIndex, treeDepth)
	if err != nil {
		return types.Digest{}, err
	}

	// We start climbing from the leaf.
	currentNode := leaf
	// When distanceFromLeaf equals treeDepth, currentNode will contain the computed root.
	for distanceFromLeaf := uint64(0); distanceFromLeaf < treeDepth; distanceFromLeaf++ {
		siblingIndexInProof := distanceFromLeaf * nodeHashSize
		// siblingHash is the next node to append to our current node, retrieved from the proof.
		siblingHash := proof[siblingIndexInProof : siblingIndexInProof+nodeHashSize]

		// Vector commitment nodes are of the form Sha256("MA" || left child || right child). To calculate the internal node,
		// we have to use the positions array to determine if our current node is the left or right child.
		// Positions[distanceFromLeaf] is the position of the current node at height distanceFromLeaf.
		nodeDomainSeparator := []byte(transactionverificationtypes.MerkleArrayNode)
		internalNodeData := nodeDomainSeparator
		switch positions[distanceFromLeaf] {
		case leftChild:
			internalNodeData = append(internalNodeData, currentNode[:]...)
			internalNodeData = append(internalNodeData, siblingHash[:]...)
		case rightChild:
			internalNodeData = append(internalNodeData, siblingHash[:]...)
			internalNodeData = append(internalNodeData, currentNode[:]...)
		default:
			return types.Digest{}, ErrInvalidPosition
		}

		currentNode = sha256.Sum256(internalNodeData)
	}

	return currentNode, nil
}

// VerifyTransaction receives a sha256 hashed transaction, a proof to compute the transaction's commitment, a proof
// to compute the commitment belonging to the light block header associated with the transaction's commitment,
// and an expected commitment to compare to. The function verifies that the computed commitment using the given proofs
// is identical to the provided commitment.
// Parameters:
// transactionHash - the result of invoking Sha256 on the msgpacked transaction.
// transactionProofResponse - the response returned by an Algorand node when queried using GetTransactionProof.
// lightBlockHeaderProofResponse - the response returned by an Algorand node when queried using the GetLightBlockHeaderProof.
// confirmedRound - the round in which the given transaction was confirmed.
// genesisHash - the hash of the genesis block.
// seed - the sortition seed of the block associated with the light block header.
// blockIntervalCommitment - the commitment to compare to, provided by the Oracle.
func VerifyTransaction(transactionHash types.Digest, transactionProofResponse models.ProofResponse,
	lightBlockHeaderProofResponse models.LightBlockHeaderProof, confirmedRound types.Round, genesisHash types.Digest, seed transactionverificationtypes.Seed, blockIntervalCommitment types.Digest) error {
	// Verifying attested vector commitment roots is currently exclusively supported with sha256 hashing, both for transactions
	// and light block headers.
	if transactionProofResponse.Hashtype != "sha256" {
		return ErrUnsupportedHashFunction
	}

	var stibHashDigest types.Digest
	copy(stibHashDigest[:], transactionProofResponse.Stibhash[:])

	// We first compute the leaf in the vector commitment that attests to the given transaction.
	transactionLeaf := computeTransactionLeaf(transactionHash, stibHashDigest)
	// We use the transactionLeaf and the given transactionProofResponse to compute the root of the vector commitment
	// that attests to the given transaction.
	transactionProofRoot, err := computeVectorCommitmentRoot(transactionLeaf, transactionProofResponse.Idx,
		transactionProofResponse.Proof, transactionProofResponse.Treedepth)

	if err != nil {
		return err
	}

	// We use our computed transaction vector commitment root, saved in transactionProofRoot, and the given data
	// to calculate the leaf in the vector commitment that attests to the light block headers.
	candidateLightBlockHeaderLeaf := computeLightBlockHeaderLeaf(confirmedRound, transactionProofRoot, genesisHash, seed)
	// We use the candidateLightBlockHeaderLeaf and the given lightBlockHeaderProofResponse to compute the root of the vector
	// commitment that attests to the candidateLightBlockHeaderLeaf.
	lightBlockHeaderProofRoot, err := computeVectorCommitmentRoot(candidateLightBlockHeaderLeaf, lightBlockHeaderProofResponse.Index, lightBlockHeaderProofResponse.Proof,
		lightBlockHeaderProofResponse.Treedepth)

	if err != nil {
		return err
	}

	// We verify that the given commitment, provided by the Oracle, is identical to the computed commitment
	if bytes.Equal(lightBlockHeaderProofRoot[:], blockIntervalCommitment[:]) != true {
		return ErrRootMismatch
	}
	return nil
}
