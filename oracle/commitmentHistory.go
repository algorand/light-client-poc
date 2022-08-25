package oracle

import (
	"errors"
	"github.com/algorand/go-stateproof-verification/types"
)

var (
	ErrTooEarlyRoundRequested = errors.New("round is earlier than the first attested round")
	ErrNoStateProofForRound   = errors.New("round belongs to an interval without a matching state proof")
)

// CommitmentHistory is our implementation of the sliding window in charge of holding block interval commitments for
// each interval, and for retrieving the appropriate commitment for a given round by calculating the interval that covers
// the given round. This implementation might look significantly different in an actual light client, as a result of
// available resources and of the environment the light client is developed in.
type CommitmentHistory struct {
	// FirstAttestedRound is the first round to which a state proof message attests.
	FirstAttestedRound uint64
	// IntervalSize is The number of rounds each state proof message attests to.
	IntervalSize uint64
	// Capacity is the maximum number of commitments to hold before discarding the earliest commitment.
	Capacity uint64
	// EarliestInterval is the earliest interval currently saved in history.
	EarliestInterval uint64
	// NextInterval is the interval to which the next state proof attest.
	NextInterval uint64
	// Data is a map of intervals to their block interval commitment.
	Data map[uint64]types.Digest
}

// InitializeCommitmentHistory initializes the commitment history using appropriate data regarding state proofs and a
// given capacity.
// Parameters:
// firstAttestedRound - the first round to which a state proof message attests.
// intervalSize - the number of rounds each state proof message attests to.
// capacity - the maximum number of commitments to hold before discarding the earliest commitment.
func InitializeCommitmentHistory(firstAttestedRound uint64, intervalSize uint64, capacity uint64) *CommitmentHistory {
	return &CommitmentHistory{
		FirstAttestedRound: firstAttestedRound,
		IntervalSize:       intervalSize,
		Capacity:           capacity,
		EarliestInterval:   0,
		NextInterval:       0,
		Data:               make(map[uint64]types.Digest),
	}
}

// GetCommitment receives a round and returns the block interval commitment for the interval that covers the given round.
// Parameters:
// round - the round to return the commitment for.
func (c *CommitmentHistory) GetCommitment(round types.Round) (types.Digest, error) {
	// Rounds earlier than state proof generation beginning can not have commitments.
	if uint64(round) < c.FirstAttestedRound {
		return types.Digest{}, ErrTooEarlyRoundRequested
	}

	// coveringInterval is the interval that covers the given round.
	coveringInterval := (uint64(round) - c.FirstAttestedRound) / c.IntervalSize

	// Intervals begin in rounds that come after rounds that are divisible by IntervalSize, so if our round is divisible
	// by IntervalSize we have to adjust its interval accordingly.
	if uint64(round)%c.IntervalSize == 0 {
		coveringInterval -= 1
	}

	// If we either don't yet have a commitment for the round or we've already discarded the commitment for the round,
	// return an error.
	if coveringInterval >= c.NextInterval || coveringInterval < c.EarliestInterval {
		return types.Digest{}, ErrNoStateProofForRound
	}

	return c.Data[coveringInterval], nil
}

func (c *CommitmentHistory) InsertCommitment(commitment types.Digest) {
	// Insert the new commitment.
	c.Data[c.NextInterval] = commitment
	// Prepare for the next commitment to be inserted.
	c.NextInterval++

	// If inserting this commitment made us exceed capacity, discard the earliest commitment we have and advance
	// EarliestInterval accordingly.
	if uint64(len(c.Data)) > c.Capacity {
		delete(c.Data, c.EarliestInterval)
		c.EarliestInterval++
	}
}
