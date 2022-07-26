package encoded_assets

import (
	"path/filepath"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/stateproofs/stateprooftypes"
	"github.com/algorand/go-algorand-sdk/types"

	"github.com/almog-t/light-client-poc/utilities"
)

func GetParsedTransactionVerificationData(transactionVerificationDataPath string) (types.Digest, types.Round, stateprooftypes.Seed, types.Digest, models.ProofResponse,
	models.LightBlockHeaderProof, types.Digest, error) {
	var genesisHash types.Digest
	err := utilities.DecodeFromFile(filepath.Join(transactionVerificationDataPath, "genesis_hash.json"), &genesisHash)
	if err != nil {
		return types.Digest{}, 0, stateprooftypes.Seed{}, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	var round types.Round
	err = utilities.DecodeFromFile(filepath.Join(transactionVerificationDataPath, "round.json"), &round)
	if err != nil {
		return types.Digest{}, 0, stateprooftypes.Seed{}, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	var seed stateprooftypes.Seed
	err = utilities.DecodeFromFile(filepath.Join(transactionVerificationDataPath, "seed.json"), &seed)
	if err != nil {
		return types.Digest{}, 0, stateprooftypes.Seed{}, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	var transactionId types.Digest
	err = utilities.DecodeFromFile(filepath.Join(transactionVerificationDataPath, "transaction_id.json"), &transactionId)
	if err != nil {
		return types.Digest{}, 0, stateprooftypes.Seed{}, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	var transactionProofResponse models.ProofResponse
	err = utilities.DecodeFromFile(filepath.Join(transactionVerificationDataPath, "transaction_proof_response.json"), &transactionProofResponse)
	if err != nil {
		return types.Digest{}, 0, stateprooftypes.Seed{}, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	var lightBlockHeaderProof models.LightBlockHeaderProof
	err = utilities.DecodeFromFile(filepath.Join(transactionVerificationDataPath, "light_block_header_proof_response.json"), &lightBlockHeaderProof)
	if err != nil {
		return types.Digest{}, 0, stateprooftypes.Seed{}, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	var lightBlockHeaderCommitment types.Digest
	err = utilities.DecodeFromFile(filepath.Join(transactionVerificationDataPath, "light_block_header_commitment.json"), &lightBlockHeaderCommitment)
	if err != nil {
		return types.Digest{}, 0, stateprooftypes.Seed{}, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	return genesisHash, round, seed, transactionId, transactionProofResponse,
		lightBlockHeaderProof, lightBlockHeaderCommitment, nil
}

func GetParsedGenesisData(genesisDataPath string) (stateprooftypes.GenericDigest, uint64, error) {
	genesisVotersCommitment := stateprooftypes.GenericDigest{}
	err := utilities.DecodeFromFile(filepath.Join(genesisDataPath, "genesis_voters_commitment.json"), &genesisVotersCommitment)
	if err != nil {
		return stateprooftypes.GenericDigest{}, 0, err
	}

	genesisVotersLnProvenWeight := uint64(0)
	err = utilities.DecodeFromFile(filepath.Join(genesisDataPath, "genesis_voters_ln_proven_weight.json"), &genesisVotersLnProvenWeight)
	if err != nil {
		return stateprooftypes.GenericDigest{}, 0, err
	}

	return genesisVotersCommitment, genesisVotersLnProvenWeight, nil
}

func GetParsedStateProofAdvancmentData(stateProofVerificationDataPath string) (stateprooftypes.Message,
	*stateprooftypes.EncodedStateProof, error) {
	stateProofMessage := stateprooftypes.Message{}
	err := utilities.DecodeFromFile(filepath.Join(stateProofVerificationDataPath, "state_proof_message.json"), &stateProofMessage)
	if err != nil {
		return stateprooftypes.Message{}, nil, err
	}

	var stateProof stateprooftypes.EncodedStateProof
	err = utilities.DecodeFromFile(filepath.Join(stateProofVerificationDataPath, "state_proof.json"), &stateProof)
	if err != nil {
		return stateprooftypes.Message{}, nil, err
	}

	return stateProofMessage, &stateProof, nil
}
