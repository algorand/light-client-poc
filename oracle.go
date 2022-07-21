package main

import (
	"errors"

	"github.com/algorand/go-algorand-sdk/stateproofs/stateprooftypes"
	"github.com/algorand/go-algorand-sdk/stateproofs/stateproofverification"
	"github.com/algorand/go-algorand-sdk/types"
)

var (
	ErrNoStateProofForRound = errors.New("round belongs to an interval without a matching state proof")
)

type Oracle struct {
	intervalSize       uint64
	firstAttestedRound uint64

	intervalCommitmentHistory []types.Digest

	stateProofVerifier *stateproofverification.StateProofVerifier
}

// InitializeOracle initializes the Oracle using trusted genesis data - the voters commitment and the Ln of the proven weight.
// These parameters can be retrieved using the SDK API.
func InitializeOracle(intervalSize uint64, firstAttestedRound uint64, genesisVotersCommitment stateprooftypes.GenericDigest, genesisLnProvenWeight uint64) *Oracle {
	stateProofVerifier := stateproofverification.InitializeVerifier(genesisVotersCommitment, genesisLnProvenWeight)

	return &Oracle{
		intervalSize:       intervalSize,
		firstAttestedRound: firstAttestedRound,

		intervalCommitmentHistory: make([]types.Digest, 0),

		stateProofVerifier: stateProofVerifier,
	}
}

func (o *Oracle) roundToInterval(round types.Round) uint64 {
	nearestIntervalMultiple := (uint64(round) / o.intervalSize) * o.intervalSize
	return (nearestIntervalMultiple - (o.firstAttestedRound - 1)) / o.intervalSize
}

// AdvanceState receives a message packed state proof, provided by the SDK API, and a state proof message that the
// state proof attests to. It verifies the message using the proof and the verifier from the previous round,
// and, if successful, creates a new verifier using the message and saves the block header commitment to the history.
func (o *Oracle) AdvanceState(stateProof *stateprooftypes.EncodedStateProof, message stateprooftypes.Message) error {
	err := o.stateProofVerifier.AdvanceState(stateProof, message)
	if err != nil {
		return err
	}

	var commitmentDigest types.Digest
	copy(commitmentDigest[:], message.BlockHeadersCommitment)
	o.intervalCommitmentHistory = append(o.intervalCommitmentHistory, commitmentDigest)

	return nil
}

func (o *Oracle) GetStateProofCommitment(round types.Round) (types.Digest, error) {
	transactionCommitmentInterval := o.roundToInterval(round)
	if transactionCommitmentInterval >= uint64(len(o.intervalCommitmentHistory)) {
		return types.Digest{}, ErrNoStateProofForRound
	}

	return o.intervalCommitmentHistory[o.roundToInterval(round)], nil
}
