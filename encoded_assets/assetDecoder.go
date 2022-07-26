package encoded_assets

import (
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/stateproofs/stateprooftypes"
	"github.com/algorand/go-algorand-sdk/types"
	"path/filepath"

	"github.com/almog-t/light-client-poc/utilities"
)

type AssetDecoder struct {
	transactionVerificationDataFolder string
	stateProofVerificationDataFolder  string
}

func InitializeDataDecoder(encodedAssetsFolder string) (*AssetDecoder, error) {
	transactionVerificationDataFolder, err := filepath.Abs(filepath.Join(encodedAssetsFolder, "transaction_verification"))
	if err != nil {
		return nil, err
	}

	stateProofVerificationDataFolder, err := filepath.Abs(filepath.Join(encodedAssetsFolder, "state_proof_verification"))
	if err != nil {
		return nil, err
	}

	return &AssetDecoder{
		transactionVerificationDataFolder: transactionVerificationDataFolder,
		stateProofVerificationDataFolder:  stateProofVerificationDataFolder,
	}, nil
}

func (d *AssetDecoder) GetParsedTransactionVerificationData() (types.Digest, types.Round, stateprooftypes.Seed, types.Digest, models.ProofResponse,
	models.LightBlockHeaderProof, types.Digest, error) {
	var genesisHash types.Digest
	err := utilities.DecodeFromFile(filepath.Join(d.transactionVerificationDataFolder, "genesis_hash.json"), &genesisHash)
	if err != nil {
		return types.Digest{}, 0, stateprooftypes.Seed{}, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	var round types.Round
	err = utilities.DecodeFromFile(filepath.Join(d.transactionVerificationDataFolder, "round.json"), &round)
	if err != nil {
		return types.Digest{}, 0, stateprooftypes.Seed{}, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	var seed stateprooftypes.Seed
	err = utilities.DecodeFromFile(filepath.Join(d.transactionVerificationDataFolder, "seed.json"), &seed)
	if err != nil {
		return types.Digest{}, 0, stateprooftypes.Seed{}, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	var transactionId types.Digest
	err = utilities.DecodeFromFile(filepath.Join(d.transactionVerificationDataFolder, "transaction_id.json"), &transactionId)
	if err != nil {
		return types.Digest{}, 0, stateprooftypes.Seed{}, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	var transactionProofResponse models.ProofResponse
	err = utilities.DecodeFromFile(filepath.Join(d.transactionVerificationDataFolder, "transaction_proof_response.json"), &transactionProofResponse)
	if err != nil {
		return types.Digest{}, 0, stateprooftypes.Seed{}, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	var lightBlockHeaderProof models.LightBlockHeaderProof
	err = utilities.DecodeFromFile(filepath.Join(d.transactionVerificationDataFolder, "light_block_header_proof_response.json"), &lightBlockHeaderProof)
	if err != nil {
		return types.Digest{}, 0, stateprooftypes.Seed{}, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	var lightBlockHeaderCommitment types.Digest
	err = utilities.DecodeFromFile(filepath.Join(d.transactionVerificationDataFolder, "light_block_header_commitment.json"), &lightBlockHeaderCommitment)
	if err != nil {
		return types.Digest{}, 0, stateprooftypes.Seed{}, types.Digest{}, models.ProofResponse{}, models.LightBlockHeaderProof{}, types.Digest{}, err
	}

	return genesisHash, round, seed, transactionId, transactionProofResponse,
		lightBlockHeaderProof, lightBlockHeaderCommitment, nil
}

func (d *AssetDecoder) GetParsedStateProofAdvancmentData() (stateprooftypes.GenericDigest, uint64, stateprooftypes.Message,
	*stateprooftypes.EncodedStateProof, error) {
	genesisVotersCommitment := stateprooftypes.GenericDigest{}
	err := utilities.DecodeFromFile(filepath.Join(d.stateProofVerificationDataFolder, "genesis_voters_commitment.json"), &genesisVotersCommitment)
	if err != nil {
		return nil, 0, stateprooftypes.Message{}, nil, err
	}

	genesisVotersLnProvenWeight := uint64(0)
	err = utilities.DecodeFromFile(filepath.Join(d.stateProofVerificationDataFolder, "genesis_voters_ln_proven_weight.json"), &genesisVotersLnProvenWeight)
	if err != nil {
		return nil, 0, stateprooftypes.Message{}, nil, err
	}

	stateProofMessage := stateprooftypes.Message{}
	err = utilities.DecodeFromFile(filepath.Join(d.stateProofVerificationDataFolder, "state_proof_message.json"), &stateProofMessage)
	if err != nil {
		return nil, 0, stateprooftypes.Message{}, nil, err
	}

	var stateProof stateprooftypes.EncodedStateProof
	err = utilities.DecodeFromFile(filepath.Join(d.stateProofVerificationDataFolder, "state_proof.json"), &stateProof)
	if err != nil {
		return nil, 0, stateprooftypes.Message{}, nil, err
	}

	return genesisVotersCommitment, genesisVotersLnProvenWeight, stateProofMessage, &stateProof, nil
}
