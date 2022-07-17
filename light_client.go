package main

import (
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/stateproofs/datatypes"
	"github.com/algorand/go-algorand-sdk/stateproofs/functions"
	"github.com/algorand/go-algorand/crypto/stateproof"
	"light-client-poc/encoded_assets"
)

const strengthTarget = uint64(256)

type LightClient struct {
	intervalSize              uint64
	firstAttestedRound        uint64
	intervalCommitmentHistory map[uint64]datatypes.GenericDigest
	genesisHash               datatypes.GenericDigest

	transactionVerifier *TransactionVerifier
	// TODO: handle import
	stateProofVerifier *stateproof.Verifier
}

func InitializeLightClient(intervalSize uint64, firstAttestedRound uint64, genesisHash datatypes.GenericDigest, genesisVotersCommitment datatypes.GenericDigest, genesisLnProvenWeight uint64) *LightClient {
	transactionVerifier := TransactionVerifier{genesisHash: genesisHash}
	stateProofVerifier := functions.MkVerifierWithLnProvenWeight(genesisVotersCommitment, genesisLnProvenWeight, strengthTarget)
	return &LightClient{
		intervalSize:              intervalSize,
		firstAttestedRound:        firstAttestedRound,
		intervalCommitmentHistory: make(map[uint64]datatypes.GenericDigest, 0),
		transactionVerifier:       &transactionVerifier,
		stateProofVerifier:        stateProofVerifier,
	}
}

func (t *LightClient) roundToInterval(round uint64) uint64 {
	nearestIntervalMultiple := (round / t.intervalSize) * t.intervalSize
	return (nearestIntervalMultiple - t.firstAttestedRound) / t.intervalSize
}

func (t *LightClient) AdvanceState(s *stateproof.StateProof, message *datatypes.Message) bool {
	messageHash := message.IntoStateProofMessageHash()

	err := t.stateProofVerifier.Verify(message.LastAttestedRound, stateproof.MessageHash(messageHash), s)
	if err != nil {
		return false
	}

	t.stateProofVerifier = functions.MkVerifierWithLnProvenWeight(message.VotersCommitment, message.LnProvenWeight, strengthTarget)
	return true
}

//TODO: Should eventually simply take the transaction itself here (after the SDK is updated)

func (t *LightClient) VerifyTransaction(transactionId datatypes.GenericDigest, round uint64, transactionProofResponse models.ProofResponse,
	lightBlockHeaderProofResponse encoded_assets.LightBlockHeaderProofResponse) (bool, error) {
	matchingCommitment := t.intervalCommitmentHistory[t.roundToInterval(round)]
	return t.transactionVerifier.VerifyTransaction(transactionId, transactionProofResponse, lightBlockHeaderProofResponse, matchingCommitment, round)
}
