package main

import (
	"errors"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/stateproofs/datatypes"
	"github.com/algorand/go-algorand-sdk/stateproofs/functions"
	"github.com/algorand/go-algorand-sdk/types"
)

var (
	ErrNoStateProofForRound = errors.New("round belongs to an interval without a matching state proof")
)

type LightClient struct {
	intervalSize       uint64
	firstAttestedRound uint64

	nextInterval              uint64
	intervalCommitmentHistory map[uint64]types.Digest

	stateProofVerifier  *functions.StateProofVerifier
	transactionVerifier *TransactionVerifier
}

func InitializeLightClient(intervalSize uint64, firstAttestedRound uint64, genesisHash types.Digest, genesisVotersCommitment datatypes.GenericDigest, genesisLnProvenWeight uint64) *LightClient {
	transactionVerifier := TransactionVerifier{genesisHash: genesisHash}
	stateProofVerifier := functions.InitializeVerifier(genesisVotersCommitment, genesisLnProvenWeight)

	return &LightClient{
		intervalSize:       intervalSize,
		firstAttestedRound: firstAttestedRound,

		nextInterval:              0,
		intervalCommitmentHistory: make(map[uint64]types.Digest, 0),

		transactionVerifier: &transactionVerifier,
		stateProofVerifier:  stateProofVerifier,
	}
}

func (l *LightClient) roundToInterval(round types.Round) uint64 {
	nearestIntervalMultiple := (uint64(round) / l.intervalSize) * l.intervalSize
	return (nearestIntervalMultiple - (l.firstAttestedRound - 1)) / l.intervalSize
}

func (l *LightClient) AdvanceState(stateProof *datatypes.EncodedStateProof, message datatypes.Message) error {
	err := l.stateProofVerifier.AdvanceState(stateProof, message)
	if err != nil {
		return err
	}

	var commitmentDigest types.Digest
	copy(commitmentDigest[:], message.BlockHeadersCommitment[:])
	l.intervalCommitmentHistory[l.nextInterval] = commitmentDigest
	l.nextInterval++

	return nil
}

func (l *LightClient) VerifyTransaction(transactionId types.Digest, transactionProofResponse models.ProofResponse,
	lightBlockHeaderProofResponse models.LightBlockHeaderProof, round types.Round, seed datatypes.Seed) error {
	matchingCommitment, ok := l.intervalCommitmentHistory[l.roundToInterval(round)]

	if !ok {
		return ErrNoStateProofForRound
	}

	return l.transactionVerifier.VerifyTransaction(transactionId, transactionProofResponse,
		lightBlockHeaderProofResponse, matchingCommitment, round, seed)
}
