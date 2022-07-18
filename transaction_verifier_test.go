package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"light-client-poc/encoded_assets"
)

func TestTransactionVerifier_VerifyTransaction(t *testing.T) {
	r := require.New(t)

	genesisHash, err := encoded_assets.GetGenesisHash()
	r.NoError(err)

	round, transactionId, transactionProofResponse, lightBlockHeaderProofResponse, lightBlockHeaderCommitment, err := encoded_assets.GetParsedTransactionVerificationData()
	r.NoError(err)

	transactionVerifier := TransactionVerifier{genesisHash: genesisHash}
	err = transactionVerifier.VerifyTransaction(transactionId, transactionProofResponse,
		lightBlockHeaderProofResponse, lightBlockHeaderCommitment, round)
	r.NoError(err)
}
