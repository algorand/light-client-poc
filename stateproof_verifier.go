package main

import (
	"github.com/algorand/go-algorand-sdk/stateproofs/datatypes"
	"github.com/algorand/go-algorand-sdk/stateproofs/functions"
	"github.com/algorand/go-algorand/crypto/stateproof"
)

const strengthTarget = uint64(256)

// TODO: rename
type StateProofTracker struct {
	intervalSize              uint64
	firstAttestedRound        uint64
	intervalCommitmentHistory map[uint64]datatypes.GenericDigest

	// TODO: handle import
	verifier *stateproof.Verifier
}

func InitializeStateProofTracker(intervalSize uint64, firstAttestedRound uint64, genesisVotersCommitment datatypes.GenericDigest, genesisLnProvenWeight uint64) *StateProofTracker {
	verifier := functions.MkVerifierWithLnProvenWeight(genesisVotersCommitment, genesisLnProvenWeight, strengthTarget)
	return &StateProofTracker{
		intervalSize:              intervalSize,
		firstAttestedRound:        firstAttestedRound,
		intervalCommitmentHistory: make(map[uint64]datatypes.GenericDigest, 0),
		verifier:                  verifier,
	}
}

func (t *StateProofTracker) roundToInterval(round uint64) uint64 {
	nearestIntervalMultiple := (round / t.intervalSize) * t.intervalSize
	return (nearestIntervalMultiple - t.firstAttestedRound) / t.intervalSize
}

func (t *StateProofTracker) ProcessStateProofAndMessage(message *datatypes.Message, round uint64, s *stateproof.StateProof) bool {
	messageHash := message.IntoStateProofMessageHash()

	err := t.verifier.Verify(round, stateproof.MessageHash(messageHash), s)
	if err != nil {
		return false
	}

	t.verifier = functions.MkVerifierWithLnProvenWeight(message.VotersCommitment, message.LnProvenWeight, strengthTarget)
	return true
}
