package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"light-client-poc/encoded_assets"
)

func getAdvancedLightClient(r *require.Assertions) *LightClient {
	genesisHash, err := encoded_assets.GetGenesisHash()

	genesisVotersCommitment, genesisVotersLnProvenWeight, stateProofMessage, stateProof, err := encoded_assets.GetParsedStateProofAdvancmentData()
	r.NoError(err)

	intervalSize := stateProofMessage.LastAttestedRound - stateProofMessage.FirstAttestedRound + 1
	tracker := InitializeLightClient(intervalSize, stateProofMessage.FirstAttestedRound, *genesisHash, *genesisVotersCommitment, genesisVotersLnProvenWeight)
	tracker.AdvanceState(stateProof, stateProofMessage)
	return tracker
}

func TestLightClient_AdvanceState(t *testing.T) {
	r := require.New(t)

	_ = getAdvancedLightClient(r)
}

func TestLightClient_VerifyTransaction(t *testing.T) {
	r := require.New(t)

	lightClient := getAdvancedLightClient(r)

	round, transactionId, transactionProofResponse, lightBlockHeaderProofResponse, _, err := encoded_assets.GetParsedTransactionVerificationData()
	r.NoError(err)

	verified, err := lightClient.VerifyTransaction(*transactionId, round, *transactionProofResponse, *lightBlockHeaderProofResponse)
	r.NoError(err)
	r.True(verified)
}
