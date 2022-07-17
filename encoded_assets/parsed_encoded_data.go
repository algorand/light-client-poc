package encoded_assets

import (
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/stateproofs/datatypes"
	"github.com/algorand/go-algorand/crypto/stateproof"
	"light-client-poc/utilities"
)

var transactionVerificationFolder = "encoded_assets/transaction_verification/"
var stateProofVerificationFolder = "encoded_assets/state_proof_verification/"
var generalAssetsFolder = "encoded_assets/general_assets/"

func GetGenesisHash() (*[]byte, error) {
	var genesisHash []byte
	err := utilities.DecodeFromFile(generalAssetsFolder+"genesis_hash.json", &genesisHash)
	if err != nil {
		return nil, err
	}

	return &genesisHash, nil
}

func GetParsedTransactionVerificationData() (uint64, *[]byte, *models.ProofResponse,
	*LightBlockHeaderProofResponse, *[]byte, error) {
	var round uint64
	err := utilities.DecodeFromFile(transactionVerificationFolder+"round.json", &round)
	if err != nil {
		return 0, nil, nil, nil, nil, err
	}

	var transactionId []byte
	err = utilities.DecodeFromFile(transactionVerificationFolder+"transaction_id.json", &transactionId)
	if err != nil {
		return 0, nil, nil, nil, nil, err
	}

	var transactionProofResponse models.ProofResponse
	err = utilities.DecodeFromFile(transactionVerificationFolder+"transaction_proof_response.json", &transactionProofResponse)
	if err != nil {
		return 0, nil, nil, nil, nil, err
	}

	var lightBlockHeaderProofResponse LightBlockHeaderProofResponse
	err = utilities.DecodeFromFile(transactionVerificationFolder+"light_block_header_proof_response.json", &lightBlockHeaderProofResponse)
	if err != nil {
		return 0, nil, nil, nil, nil, err
	}

	var lightBlockHeaderCommitment []byte
	err = utilities.DecodeFromFile(transactionVerificationFolder+"light_block_header_commitment.json", &lightBlockHeaderCommitment)
	if err != nil {
		return 0, nil, nil, nil, nil, err
	}

	return round, &transactionId, &transactionProofResponse,
		&lightBlockHeaderProofResponse, &lightBlockHeaderCommitment, nil
}

func GetParsedStateProofAdvancmentData() (*datatypes.GenericDigest, uint64, *datatypes.Message,
	*stateproof.StateProof, error) {
	genesisVotersCommitment := datatypes.GenericDigest{}
	err := utilities.DecodeFromFile(stateProofVerificationFolder+"genesis_voters_commitment.json", &genesisVotersCommitment)
	if err != nil {
		return nil, 0, nil, nil, err
	}

	genesisVotersLnProvenWeight := uint64(0)
	err = utilities.DecodeFromFile(stateProofVerificationFolder+"genesis_voters_ln_proven_weight.json", &genesisVotersLnProvenWeight)
	if err != nil {
		return nil, 0, nil, nil, err
	}

	stateProofMessage := datatypes.Message{}
	err = utilities.DecodeFromFile(stateProofVerificationFolder+"state_proof_message.json", &stateProofMessage)
	if err != nil {
		return nil, 0, nil, nil, err
	}

	stateProof := stateproof.StateProof{}
	err = utilities.DecodeFromFile(stateProofVerificationFolder+"state_proof.json", &stateProof)
	if err != nil {
		return nil, 0, nil, nil, err
	}

	return &genesisVotersCommitment, genesisVotersLnProvenWeight, &stateProofMessage, &stateProof, nil
}
