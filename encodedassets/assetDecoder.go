package encodedassets

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/types"

	"github.com/algorand/go-stateproof-verification/stateproof"
	"github.com/algorand/go-stateproof-verification/stateproofcrypto"
)

// These functions take the encoded assets, committed as examples, and parse them.

func decodeFromFile(encodedPath string, target interface{}) error {
	encodedData, err := os.ReadFile(encodedPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(encodedData, target)
	return err
}

func GetParsedtypesData(typesDataPath string) (types.Digest, types.Round, types.Seed, types.Digest, models.TransactionProofResponse,
	models.LightBlockHeaderProof, error) {
	var genesisHash types.Digest
	err := decodeFromFile(filepath.Join(typesDataPath, "genesis_hash.txt"), &genesisHash)
	if err != nil {
		return types.Digest{}, 0, types.Seed{}, types.Digest{}, models.TransactionProofResponse{}, models.LightBlockHeaderProof{}, err
	}

	var round types.Round
	err = decodeFromFile(filepath.Join(typesDataPath, "round.txt"), &round)
	if err != nil {
		return types.Digest{}, 0, types.Seed{}, types.Digest{}, models.TransactionProofResponse{}, models.LightBlockHeaderProof{}, err
	}

	var seed types.Seed
	err = decodeFromFile(filepath.Join(typesDataPath, "seed.txt"), &seed)
	if err != nil {
		return types.Digest{}, 0, types.Seed{}, types.Digest{}, models.TransactionProofResponse{}, models.LightBlockHeaderProof{}, err
	}

	var transactionId types.Digest
	err = decodeFromFile(filepath.Join(typesDataPath, "transaction_id.txt"), &transactionId)
	if err != nil {
		return types.Digest{}, 0, types.Seed{}, types.Digest{}, models.TransactionProofResponse{}, models.LightBlockHeaderProof{}, err
	}

	var transactionProofResponse models.TransactionProofResponse
	err = decodeFromFile(filepath.Join(typesDataPath, "transaction_proof_response.json"), &transactionProofResponse)
	if err != nil {
		return types.Digest{}, 0, types.Seed{}, types.Digest{}, models.TransactionProofResponse{}, models.LightBlockHeaderProof{}, err
	}

	var lightBlockHeaderProof models.LightBlockHeaderProof
	err = decodeFromFile(filepath.Join(typesDataPath, "light_block_header_proof_response.json"), &lightBlockHeaderProof)
	if err != nil {
		return types.Digest{}, 0, types.Seed{}, types.Digest{}, models.TransactionProofResponse{}, models.LightBlockHeaderProof{}, err
	}

	return genesisHash, round, seed, transactionId, transactionProofResponse,
		lightBlockHeaderProof, nil
}

func GetParsedGenesisData(genesisDataPath string) (stateproofcrypto.GenericDigest, uint64, error) {
	genesisVotersCommitment := stateproofcrypto.GenericDigest{}
	err := decodeFromFile(filepath.Join(genesisDataPath, "genesis_voters_commitment.txt"), &genesisVotersCommitment)
	if err != nil {
		return stateproofcrypto.GenericDigest{}, 0, err
	}

	genesisVotersLnProvenWeight := uint64(0)
	err = decodeFromFile(filepath.Join(genesisDataPath, "genesis_voters_ln_proven_weight.txt"), &genesisVotersLnProvenWeight)
	if err != nil {
		return stateproofcrypto.GenericDigest{}, 0, err
	}

	return genesisVotersCommitment, genesisVotersLnProvenWeight, nil
}

func GetParsedStateProofAdvancmentData(stateProofVerificationDataPath string) (types.Message,
	*stateproof.StateProof, error) {
	stateProofMessage := types.Message{}
	err := decodeFromFile(filepath.Join(stateProofVerificationDataPath, "state_proof_message.json"), &stateProofMessage)
	if err != nil {
		return types.Message{}, nil, err
	}

	var msgPackedStateProof []byte
	err = decodeFromFile(filepath.Join(stateProofVerificationDataPath, "state_proof.txt"), &msgPackedStateProof)
	if err != nil {
		return types.Message{}, nil, err
	}

	var stateProof stateproof.StateProof
	err = msgpack.Decode(msgPackedStateProof, &stateProof)
	if err != nil {
		return types.Message{}, nil, err
	}

	return stateProofMessage, &stateProof, nil
}
