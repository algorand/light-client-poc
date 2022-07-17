package main

import (
	"testing"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/stretchr/testify/require"
)

var transactionVerificationFolder = "encoded_assets/transaction_verification/"

func TestTransactionVerify(t *testing.T) {
	r := require.New(t)

	var genesisHash []byte
	err := decodeFromFile(transactionVerificationFolder+"genesis_hash.json", &genesisHash)
	r.NoError(err)

	transactionVerifier := TransactionVerifier{genesisHash: genesisHash}

	var round uint64
	err = decodeFromFile(transactionVerificationFolder+"round.json", &round)
	r.NoError(err)

	var transactionId []byte
	err = decodeFromFile(transactionVerificationFolder+"transaction_id.json", &transactionId)
	r.NoError(err)

	var transactionProofResponse models.ProofResponse
	err = decodeFromFile(transactionVerificationFolder+"transaction_proof_response.json", &transactionProofResponse)
	r.NoError(err)

	var lightBlockHeaderProofResponse LightBlockHeaderProofResponse
	err = decodeFromFile(transactionVerificationFolder+"light_block_header_proof_response.json", &lightBlockHeaderProofResponse)
	r.NoError(err)

	// TODO: Get this from the state proof trusted data
	var lightBlockHeaderCommitment []byte
	err = decodeFromFile(transactionVerificationFolder+"light_block_header_commitment.json", &lightBlockHeaderCommitment)
	r.NoError(err)

	verified, err := transactionVerifier.VerifyTransaction(transactionId, transactionProofResponse,
		lightBlockHeaderProofResponse, lightBlockHeaderCommitment, round)
	r.NoError(err)
	if !verified {
		t.Fail()
	}
}
