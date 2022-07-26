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
	BlockIntervalCommitmentHistory *utilities.CommitmentHistory
	VotersCommitment               stateprooftypes.GenericDigest
	LnProvenWeight                 uint64
}

// InitializeOracle initializes the Oracle using trusted genesis data - the voters commitment and the Ln of the proven weight.
// These parameters can be found in the developer portal.
func InitializeOracle(intervalSize uint64, genesisVotersCommitment stateprooftypes.GenericDigest,
	genesisLnProvenWeight uint64, capacity uint64) *Oracle {

	return &Oracle{
		BlockIntervalCommitmentHistory: utilities.InitializeCommitmentHistory(intervalSize, capacity),
		VotersCommitment:               genesisVotersCommitment,
		LnProvenWeight:                 genesisLnProvenWeight,
	}
}

// AdvanceState receives a message packed state proof, provided by the SDK API, and a state proof message that the
// state proof attests to. It verifies the message using the proof and the verifier from the previous round,
// and, if successful, updates the Oracle's state using the message and saves the block header commitment to the history.
func (o *Oracle) AdvanceState(stateProof *stateprooftypes.EncodedStateProof, message stateprooftypes.Message) error {
	verifier := stateproofverification.InitializeVerifier(o.VotersCommitment, o.LnProvenWeight)
	err := verifier.VerifyStateProofMessage(stateProof, message)
	if err != nil {
		return err
	}

	o.VotersCommitment = message.VotersCommitment
	o.LnProvenWeight = message.LnProvenWeight

	var commitmentDigest types.Digest
	copy(commitmentDigest[:], message.BlockHeadersCommitment)
	o.BlockIntervalCommitmentHistory.InsertCommitment(commitmentDigest)

	return nil
}

func (o *Oracle) GetStateProofCommitment(round types.Round) (types.Digest, error) {
	return o.BlockIntervalCommitmentHistory.GetCommitment(round)
}
