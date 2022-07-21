package main

import (
	"errors"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/stateproofs/stateprooftypes"
	"github.com/algorand/go-algorand-sdk/stateproofs/stateproofverification"
	"github.com/algorand/go-algorand-sdk/types"
)

var (
	ErrNoStateProofForRound = errors.New("round belongs to an interval without a matching state proof")
)

type LightClient struct {
	intervalSize       uint64
	firstAttestedRound uint64

	intervalCommitmentHistory []types.Digest

	stateProofVerifier  *stateproofverification.StateProofVerifier
	transactionVerifier *TransactionVerifier
}

// TODO: What is the parameter, where does it come from?
func InitializeLightClient(intervalSize uint64, firstAttestedRound uint64, genesisHash types.Digest, genesisVotersCommitment stateprooftypes.GenericDigest, genesisLnProvenWeight uint64) *LightClient {
	transactionVerifier := TransactionVerifier{genesisHash: genesisHash}
	stateProofVerifier := stateproofverification.InitializeVerifier(genesisVotersCommitment, genesisLnProvenWeight)

	return &LightClient{
		intervalSize:       intervalSize,
		firstAttestedRound: firstAttestedRound,

		intervalCommitmentHistory: make([]types.Digest, 0),

		transactionVerifier: &transactionVerifier,
		stateProofVerifier:  stateProofVerifier,
	}
}

func (l *LightClient) roundToInterval(round types.Round) uint64 {
	nearestIntervalMultiple := (uint64(round) / l.intervalSize) * l.intervalSize
	return (nearestIntervalMultiple - (l.firstAttestedRound - 1)) / l.intervalSize
}

// TODO: What is the parameter, where does it come from?
func (l *LightClient) AdvanceState(stateProof *stateprooftypes.EncodedStateProof, message stateprooftypes.Message) error {
	err := l.stateProofVerifier.AdvanceState(stateProof, message)
	if err != nil {
		return err
	}

	var commitmentDigest types.Digest
	copy(commitmentDigest[:], message.BlockHeadersCommitment)
	l.intervalCommitmentHistory = append(l.intervalCommitmentHistory, commitmentDigest)

	return nil
}

// TODO: What is the parameter, where does it come from?
func (l *LightClient) VerifyTransaction(transactionId types.Digest, transactionProofResponse models.ProofResponse,
	lightBlockHeaderProofResponse models.LightBlockHeaderProof, round types.Round, seed stateprooftypes.Seed) error {
	transactionCommitmentInterval := l.roundToInterval(round)
	if transactionCommitmentInterval >= uint64(len(l.intervalCommitmentHistory)) {
		return ErrNoStateProofForRound
	}

	matchingCommitment := l.intervalCommitmentHistory[l.roundToInterval(round)]

	return l.transactionVerifier.VerifyTransaction(transactionId, transactionProofResponse, lightBlockHeaderProofResponse, round, seed, matchingCommitment)
}
