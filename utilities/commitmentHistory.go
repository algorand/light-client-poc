package utilities

import (
	"errors"
	"github.com/algorand/go-algorand-sdk/types"
)

var (
	ErrNoStateProofForRound = errors.New("round belongs to an interval without a matching state proof")
)

type CommitmentHistory struct {
	IntervalSize     uint64
	Capacity         uint64
	EarliestInterval uint64
	NextInterval     uint64
	Data             map[uint64]types.Digest
}

// InitializeCommitmentHistory initializes the commitment history with the interval size and the capacity. Note that
// in an actual light client, some calculation would have to be made to initialize the earliest interval correctly.
func InitializeCommitmentHistory(intervalSize uint64, capacity uint64) *CommitmentHistory {
	return &CommitmentHistory{
		IntervalSize:     intervalSize,
		Capacity:         capacity,
		EarliestInterval: 0,
		NextInterval:     0,
		Data:             make(map[uint64]types.Digest),
	}
}

func (c *CommitmentHistory) GetCommitment(round types.Round) (types.Digest, error) {
	nearestInterval := (uint64(round) / c.IntervalSize) - 1
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
