package utilities

import (
	"errors"
	"github.com/algorand/go-algorand-sdk/types"
)

var (
	ErrNoStateProofForRound = errors.New("round belongs to an interval without a matching state proof")
)

type CommitmentHistory struct {
	intervalSize     uint64
	capacity         uint64
	earliestInterval uint64
	nextInterval     uint64
	data             map[uint64]types.Digest
}

// InitializeCommitmentHistory initializes the commitment history with the interval size and the capacity. Note that
// in an actual light client, some calculation would have to be made to initialize the earliest interval correctly.
func InitializeCommitmentHistory(intervalSize uint64, capacity uint64) *CommitmentHistory {
	return &CommitmentHistory{
		intervalSize:     intervalSize,
		capacity:         capacity,
		earliestInterval: 0,
		nextInterval:     0,
		data:             make(map[uint64]types.Digest),
	}
}

func (c *CommitmentHistory) GetCommitment(round types.Round) (types.Digest, error) {
	nearestInterval := (uint64(round) / c.intervalSize) - 1
	if nearestInterval >= c.nextInterval || nearestInterval < c.earliestInterval {
		return types.Digest{}, ErrNoStateProofForRound
	}
	return c.data[nearestInterval], nil
}

func (c *CommitmentHistory) InsertCommitment(commitment types.Digest) {
	c.data[c.nextInterval] = commitment
	c.nextInterval++
	if uint64(len(c.data)) > c.capacity {
		delete(c.data, c.earliestInterval)
		c.earliestInterval++
	}
}
