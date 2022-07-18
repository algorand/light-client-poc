package encoded_assets

import (
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/stateproofs/datatypes"
	"github.com/algorand/go-algorand-sdk/types"
	"light-client-poc/utilities"
)

var transactionVerificationFolder = "encoded_assets/transaction_verification/"
var stateProofVerificationFolder = "encoded_assets/state_proof_verification/"
var generalAssetsFolder = "encoded_assets/general_assets/"

func GetGenesisHash() (types.Digest, error) {
	var genesisHash types.Digest
	err := utilities.DecodeFromFile(generalAssetsFolder+"genesis_hash.json", &genesisHash)
	if err != nil {
		return types.Digest{}, err
	}

	return genesisHash, nil
}

func GetParsedTransactionVerificationData() (types.Round, types.Digest, models.ProofResponse,
	models.LightBlockHeaderProof, types.Digest, error) {
	var round types.Round
	err := utilities.DecodeFromFile(transactionVerificationFolder+"round.json", &round)
	if err != nil {
		return 0, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	var transactionId types.Digest
	err = utilities.DecodeFromFile(transactionVerificationFolder+"transaction_id.json", &transactionId)
	if err != nil {
		return 0, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	var transactionProofResponse models.ProofResponse
	err = utilities.DecodeFromFile(transactionVerificationFolder+"transaction_proof_response.json", &transactionProofResponse)
	if err != nil {
		return 0, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	var lightBlockHeaderProof models.LightBlockHeaderProof
	err = utilities.DecodeFromFile(transactionVerificationFolder+"light_block_header_proof_response.json", &lightBlockHeaderProof)
	if err != nil {
		return 0, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	var lightBlockHeaderCommitment types.Digest
	err = utilities.DecodeFromFile(transactionVerificationFolder+"light_block_header_commitment.json", &lightBlockHeaderCommitment)
	if err != nil {
		return 0, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	return round, transactionId, transactionProofResponse,
		lightBlockHeaderProof, lightBlockHeaderCommitment, nil
}

func GetParsedStateProofAdvancmentData() (datatypes.GenericDigest, uint64, datatypes.Message,
	*datatypes.EncodedStateProof, error) {
	genesisVotersCommitment := datatypes.GenericDigest{}
	err := utilities.DecodeFromFile(stateProofVerificationFolder+"genesis_voters_commitment.json", &genesisVotersCommitment)
	if err != nil {
		return nil, 0, datatypes.Message{}, nil, err
	}

	genesisVotersLnProvenWeight := uint64(0)
	err = utilities.DecodeFromFile(stateProofVerificationFolder+"genesis_voters_ln_proven_weight.json", &genesisVotersLnProvenWeight)
	if err != nil {
		return nil, 0, datatypes.Message{}, nil, err
	}

	stateProofMessage := datatypes.Message{}
	err = utilities.DecodeFromFile(stateProofVerificationFolder+"state_proof_message.json", &stateProofMessage)
	if err != nil {
		return nil, 0, datatypes.Message{}, nil, err
	}

	var stateProof datatypes.EncodedStateProof
	err = utilities.DecodeFromFile(stateProofVerificationFolder+"state_proof.json", &stateProof)
	if err != nil {
		return nil, 0, datatypes.Message{}, nil, err
	}

	return genesisVotersCommitment, genesisVotersLnProvenWeight, stateProofMessage, &stateProof, nil
}
