package main

import (
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/stateproofs/datatypes"
	"github.com/algorand/go-algorand-sdk/stateproofs/functions"
	"github.com/algorand/go-algorand-sdk/types"
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

// TODO: add error in nonexistent interval
func (t *LightClient) roundToInterval(round types.Round) uint64 {
	nearestIntervalMultiple := (uint64(round) / t.intervalSize) * t.intervalSize
	return (nearestIntervalMultiple - (t.firstAttestedRound - 1)) / t.intervalSize
}

func (t *LightClient) AdvanceState(stateProof *datatypes.EncodedStateProof, message datatypes.Message) error {
	err := t.stateProofVerifier.AdvanceState(stateProof, message)
	if err != nil {
		return err
	}

	var commitmentDigest types.Digest
	copy(commitmentDigest[:], message.BlockHeadersCommitment[:])
	t.intervalCommitmentHistory[t.nextInterval] = commitmentDigest
	t.nextInterval++

	return nil
}

func (t *LightClient) VerifyTransaction(transactionId types.Digest, transactionProofResponse models.ProofResponse,
	lightBlockHeaderProofResponse models.LightBlockHeaderProof, round types.Round) error {
	matchingCommitment := t.intervalCommitmentHistory[t.roundToInterval(round)]
	return t.transactionVerifier.VerifyTransaction(transactionId, transactionProofResponse, lightBlockHeaderProofResponse, matchingCommitment, round)
}
