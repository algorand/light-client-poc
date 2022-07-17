package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algorand/go-algorand-sdk/stateproofs/datatypes"
	"github.com/algorand/go-algorand/crypto/stateproof"
)

var stateProofVerificationFolder = "encoded_assets/state_proof_verification/"

func TestStateProofVerification(t *testing.T) {
	r := require.New(t)

	genesisVotersCommitment := datatypes.GenericDigest{}
	err := decodeFromFile(stateProofVerificationFolder+"genesis_voters_commitment.json", &genesisVotersCommitment)
	r.NoError(err)

	genesisVotersLnProvenWeight := uint64(0)
	err = decodeFromFile(stateProofVerificationFolder+"genesis_voters_ln_proven_weight.json", &genesisVotersLnProvenWeight)
	r.NoError(err)

	stateProofMessage := datatypes.Message{}
	err = decodeFromFile(stateProofVerificationFolder+"state_proof_message.json", &stateProofMessage)
	r.NoError(err)

	stateProof := stateproof.StateProof{}
	err = decodeFromFile(stateProofVerificationFolder+"state_proof.json", &stateProof)
	r.NoError(err)

	tracker := InitializeStateProofTracker(8, stateProofMessage.FirstAttestedRound, genesisVotersCommitment, genesisVotersLnProvenWeight)
	tracker.ProcessStateProofAndMessage(&stateProofMessage, 16, &stateProof)
}
