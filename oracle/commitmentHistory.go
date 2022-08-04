package oracle

import (
	"errors"
	"github.com/algorand/go-algorand-sdk/types"
)

var (
	ErrNoStateProofForRound = errors.New("round belongs to an interval without a matching state proof")
)

// CommitmentHistory is our implementation of the sliding window in charge of holding block interval commitments for
// each interval, and for retrieving the appropriate commitment for a given round by calculating the interval that covers
// the given round. This implementation might look significantly different in an actual light client, as a result of
// available resources and of the environment the light client is developed in.
type CommitmentHistory struct {
	// IntervalSize is The number of rounds each state proof message attests to.
	IntervalSize uint64
	// Capacity is the capacity of the window - inserting a commitment when the window is at capacity will cause the
	// earliest commitment to be discarded.
	Capacity uint64
	// EarliestInterval is the earliest interval currently saved in history.
	EarliestInterval uint64
	// NextInterval is the interval to which the next state proof attest.
	NextInterval uint64
	// Data is a map of intervals to their block interval commitment.
	Data map[uint64]types.Digest
}

// InitializeCommitmentHistory initializes the commitment history with the interval size and the capacity.
func InitializeCommitmentHistory(intervalSize uint64, capacity uint64) *CommitmentHistory {
	return &CommitmentHistory{
		IntervalSize:     intervalSize,
		Capacity:         capacity,
		EarliestInterval: 1,
		NextInterval:     1,
		Data:             make(map[uint64]types.Digest),
	}
}

func (c *CommitmentHistory) GetCommitment(round types.Round) (types.Digest, error) {
	nearestInterval := uint64(round) / c.IntervalSize
	if uint64(round)%c.IntervalSize == 0 {
		nearestInterval -= 1
	}

	if nearestInterval >= c.NextInterval || nearestInterval < c.EarliestInterval {
		return types.Digest{}, ErrNoStateProofForRound
	}
	return c.Data[nearestInterval], nil
}

func (c *CommitmentHistory) InsertCommitment(commitment types.Digest) {
	c.Data[c.NextInterval] = commitment
	c.NextInterval++
	if uint64(len(c.Data)) > c.Capacity {
		delete(c.Data, c.EarliestInterval)
		c.EarliestInterval++
	}
}
