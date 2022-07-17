package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	
	"light-client-poc/encoded_assets"
)

func getExampleLightClient(r *require.Assertions) *LightClient {
	genesisVotersCommitment, genesisVotersLnProvenWeight, stateProofMessage, stateProof, err := encoded_assets.GetParsedStateProofAdvancmentData()
	r.NoError(err)

	intervalSize := stateProofMessage.LastAttestedRound - stateProofMessage.FirstAttestedRound + 1
	tracker := InitializeLightClient(intervalSize, stateProofMessage.FirstAttestedRound, []byte{}, *genesisVotersCommitment, genesisVotersLnProvenWeight)
	tracker.AdvanceState(stateProof, stateProofMessage)
	return tracker
}

func TestLightClient_AdvanceState(t *testing.T) {
	r := require.New(t)

	_ = getExampleLightClient(r)
}

//func TestLightClient_VerifyTransaction(t *testing.T) {
//	r := require.New(t)
//
//	lightClient := getExampleLightClient(r)
//
//	lightClient.
//}
