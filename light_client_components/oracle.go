package light_client_components

import (
	"github.com/algorand/go-algorand-sdk/stateproofs/stateprooftypes"
	"github.com/algorand/go-algorand-sdk/stateproofs/stateproofverification"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/almog-t/light-client-poc/utilities"
)

// Oracle is in charge of maintaining commitments for previous round intervals and allows, given a round, to retrieve
// the vector commitment attesting to that round.
type Oracle struct {
	// BlockIntervalCommitmentHistory is a sliding window that holds block interval commitments for each interval. Given a round,
	// it calculates that round's interval and retrieves the block interval commitment for the calculated interval.
	BlockIntervalCommitmentHistory *utilities.CommitmentHistory
	// VotersCommitment is the current voters commitment that will be used to verify the next state proof.
	VotersCommitment stateprooftypes.GenericDigest
	// LnProvenWeight is the ln of the current proven weight. This value will be used to verify the next state proof.
	LnProvenWeight uint64
}

// InitializeOracle initializes the Oracle using trusted genesis data.
// Parameters:
// intervalSize represents the number of rounds that occur between each state proof.
// Both genesisVotersCommitment and genesisLnProvenWeight can be found in the developer portal.
// capacity is the BlockIntervalCommitmentHistory window size
func InitializeOracle(intervalSize uint64, genesisVotersCommitment stateprooftypes.GenericDigest,
	genesisLnProvenWeight uint64, capacity uint64) *Oracle {
	return &Oracle{
		// The BlockIntervalCommitmentHistory is initialized using the interval size (to calculate a given round's interval
		// when retrieving that round's commitment) and its capacity.
		BlockIntervalCommitmentHistory: utilities.InitializeCommitmentHistory(intervalSize, capacity),
		// VotersCommitment is initialized using genesis data, which can be found in Algorand's developer portal.
		VotersCommitment: genesisVotersCommitment,
		// LnProvenWeight is initialized using genesis data, which can be found in Algorand's developer portal.
		LnProvenWeight: genesisLnProvenWeight,
	}
}

// AdvanceState receives a message packed state proof, provided by the SDK API, and a state proof message that the
// state proof attests to. It verifies the message using the proof and the verifier from the previous round,
// and, if successful, updates the Oracle's state using the message and saves the block header commitment to the history.
// This method should be called by a relay or some external process that is initiated when new Algorand state proofs are available.
// Parameters:
// stateProof is a slice containing the message packed state proof, as returned from the SDK API.
// message is the decoded state proof message, unpacked using msgpack.
func (o *Oracle) AdvanceState(stateProof *stateprooftypes.EncodedStateProof, message stateprooftypes.Message) error {
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
// round is the round to which a commitment will be retrieved.
func (o *Oracle) GetStateProofCommitment(round types.Round) (types.Digest, error) {
	// Receiving a commitment that should cover a round requires calculating the round's interval and retrieving the commitment
	// for that interval. See BlockIntervalCommitmentHistory.GetCommitment for more details.
	return o.BlockIntervalCommitmentHistory.GetCommitment(round)
}
