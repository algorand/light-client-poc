package oracle

import (
	"github.com/algorand/go-algorand-sdk/stateproofs/stateproofverification"
	"github.com/algorand/go-algorand-sdk/stateproofs/transactionverificationtypes"
	"github.com/algorand/go-algorand-sdk/types"
)

// Oracle is in charge of maintaining commitments for previous round intervals and allows, given a round, to retrieve
// the vector commitment root attesting to that round.
type Oracle struct {
	// BlockIntervalCommitmentHistory is a sliding window that holds block interval commitments for each interval. Given a round,
	// it calculates that round's interval and retrieves the block interval commitment for the calculated interval.
	BlockIntervalCommitmentHistory *CommitmentHistory
	// VotersCommitment is the vector commitment root of the top N accounts to sign the next StateProof.
	VotersCommitment transactionverificationtypes.GenericDigest
	// LnProvenWeight is an integer value representing the natural log of the proven weight with 16 bits of precision.
	// This value would be used to verify the next state proof.
	LnProvenWeight uint64
}

// InitializeOracle initializes the Oracle using trusted genesis data.
// Parameters:
// firstAttestedRound - the first round to which a state proof message attests.
// intervalSize - represents the number of rounds that occur between each state proof.
// genesisVotersCommitment - the initial genesisVotersCommitment commitment. Real values can be found in the Algorand developer portal.
// genesisLnProvenWeight - the initial LnProvenWeight. Real values can be found in the Algorand developer portal.
// capacity - the maximum number of commitments to hold before discarding the earliest commitment.
func InitializeOracle(firstAttestedRound uint64, intervalSize uint64, genesisVotersCommitment transactionverificationtypes.GenericDigest,
	genesisLnProvenWeight uint64, capacity uint64) *Oracle {
	return &Oracle{
		// The BlockIntervalCommitmentHistory is initialized using the first attested round,
		// the interval size and its capacity.
		BlockIntervalCommitmentHistory: InitializeCommitmentHistory(firstAttestedRound, intervalSize, capacity),
		VotersCommitment:               genesisVotersCommitment,
		LnProvenWeight:                 genesisLnProvenWeight,
	}
}

// AdvanceState receives a message packed state proof, provided by the SDK API, and a state proof message that the
// state proof attests to. It verifies the message using the proof given and the VotersCommitment and LnProvenWeight
// from the previous state proof message.
// If successful, it updates the Oracle's VotersCommitment and LnProvenWeight using their values from the new message,
// and saves the block header commitment to the history.
// This method should be called by a relay or some external process that is initiated when new Algorand state proofs are available.
// Parameters:
// stateProof - a slice containing the message packed state proof, as returned from the SDK API.
// message - the decoded state proof message, unpacked using msgpack.
func (o *Oracle) AdvanceState(stateProof *transactionverificationtypes.EncodedStateProof, message transactionverificationtypes.Message) error {
	// verifier is Algorand's implementation of the state proof verifier, exposed by the SDK. It uses the
	// previous proven VotersCommitment and LnProvenWeight.
	verifier := stateproofverification.InitializeVerifier(o.VotersCommitment, o.LnProvenWeight)
	// The newly formed verifier verifies the given message using the state proof.
	err := verifier.VerifyStateProofMessage(stateProof, message)
	if err != nil {
		// If the verification failed, for whatever reason, we return the error returned.
		return err
	}

	// Successful verification of the message means we can trust it, so we save the VotersCommitment
	// and the LnProvenWeight in the message, for verification of the next message.
	o.VotersCommitment = message.VotersCommitment
	o.LnProvenWeight = message.LnProvenWeight

	var commitmentDigest types.Digest
	copy(commitmentDigest[:], message.BlockHeadersCommitment)
	// We insert the BlockHeadersCommitment found in the message to our commitment history sliding window.
	// A side effect of this, if this commitment were to push our window over its capacity, would be deletion
	// of the earliest commitment.
	o.BlockIntervalCommitmentHistory.InsertCommitment(commitmentDigest)

	return nil
}

// GetStateProofCommitment retrieves a saved commitment for a specific round.
// Parameters:
// round - the round to which a commitment will be retrieved.
func (o *Oracle) GetStateProofCommitment(round types.Round) (types.Digest, error) {
	// Receiving a commitment that should cover a round requires calculating the round's interval and retrieving the commitment
	// for that interval. See BlockIntervalCommitmentHistory.GetCommitment for more details.
	return o.BlockIntervalCommitmentHistory.GetCommitment(round)
}
