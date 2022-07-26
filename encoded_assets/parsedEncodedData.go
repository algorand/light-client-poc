package encoded_assets

import (
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/stateproofs/stateprooftypes"
	"github.com/algorand/go-algorand-sdk/types"
	"path/filepath"

	"github.com/almog-t/light-client-poc/utilities"
)

func GetParsedTransactionVerificationData() (types.Digest, types.Round, stateprooftypes.Seed, types.Digest, models.ProofResponse,
	models.LightBlockHeaderProof, types.Digest, error) {
	transactionVerificationFolder, err := filepath.Abs(filepath.Join("encoded_assets", "transaction_verification"))
	if err != nil {
		return types.Digest{}, 0, stateprooftypes.Seed{}, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	var genesisHash types.Digest
	err = utilities.DecodeFromFile(filepath.Join(transactionVerificationFolder, "genesis_hash.json"), &genesisHash)
	if err != nil {
		return types.Digest{}, 0, stateprooftypes.Seed{}, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	var round types.Round
	err = utilities.DecodeFromFile(filepath.Join(transactionVerificationFolder, "round.json"), &round)
	if err != nil {
		return types.Digest{}, 0, stateprooftypes.Seed{}, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	var seed stateprooftypes.Seed
	err = utilities.DecodeFromFile(filepath.Join(transactionVerificationFolder, "seed.json"), &seed)
	if err != nil {
		return types.Digest{}, 0, stateprooftypes.Seed{}, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	var transactionId types.Digest
	err = utilities.DecodeFromFile(filepath.Join(transactionVerificationFolder, "transaction_id.json"), &transactionId)
	if err != nil {
		return types.Digest{}, 0, stateprooftypes.Seed{}, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	var transactionProofResponse models.ProofResponse
	err = utilities.DecodeFromFile(filepath.Join(transactionVerificationFolder, "transaction_proof_response.json"), &transactionProofResponse)
	if err != nil {
		return types.Digest{}, 0, stateprooftypes.Seed{}, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	var lightBlockHeaderProof models.LightBlockHeaderProof
	err = utilities.DecodeFromFile(filepath.Join(transactionVerificationFolder, "light_block_header_proof_response.json"), &lightBlockHeaderProof)
	if err != nil {
		return types.Digest{}, 0, stateprooftypes.Seed{}, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	var lightBlockHeaderCommitment types.Digest
	err = utilities.DecodeFromFile(filepath.Join(transactionVerificationFolder, "light_block_header_commitment.json"), &lightBlockHeaderCommitment)
	if err != nil {
		return types.Digest{}, 0, stateprooftypes.Seed{}, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	return genesisHash, round, seed, transactionId, transactionProofResponse,
		lightBlockHeaderProof, lightBlockHeaderCommitment, nil
}

func GetParsedStateProofAdvancmentData() (stateprooftypes.GenericDigest, uint64, stateprooftypes.Message,
	*stateprooftypes.EncodedStateProof, error) {
	stateProofVerificationFolder, err := filepath.Abs(filepath.Join("encoded_assets", "state_proof_verification"))
	if err != nil {
		return nil, 0, stateprooftypes.Message{}, nil, err
	}

	genesisVotersCommitment := stateprooftypes.GenericDigest{}
	err = utilities.DecodeFromFile(filepath.Join(stateProofVerificationFolder, "genesis_voters_commitment.json"), &genesisVotersCommitment)
	if err != nil {
		return nil, 0, stateprooftypes.Message{}, nil, err
	}

	genesisVotersLnProvenWeight := uint64(0)
	err = utilities.DecodeFromFile(filepath.Join(stateProofVerificationFolder, "genesis_voters_ln_proven_weight.json"), &genesisVotersLnProvenWeight)
	if err != nil {
		return nil, 0, stateprooftypes.Message{}, nil, err
	}

	stateProofMessage := stateprooftypes.Message{}
	err = utilities.DecodeFromFile(filepath.Join(stateProofVerificationFolder, "state_proof_message.json"), &stateProofMessage)
	if err != nil {
		return nil, 0, stateprooftypes.Message{}, nil, err
	}

	var stateProof stateprooftypes.EncodedStateProof
	err = utilities.DecodeFromFile(filepath.Join(stateProofVerificationFolder, "state_proof.json"), &stateProof)
	if err != nil {
		return nil, 0, stateprooftypes.Message{}, nil, err
	}

	return genesisVotersCommitment, genesisVotersLnProvenWeight, stateProofMessage, &stateProof, nil
}
