package main

import (
	"fmt"
	"github.com/almog-t/light-client-poc/encoded_assets"
	"github.com/almog-t/light-client-poc/oracle"
	"github.com/almog-t/light-client-poc/transactionverifier"
)

// A light client is composed of two modules:
// 1. An oracle, in charge of maintaining Algorand's state as verified with state proofs. For more details, see oracle.go.
// 2. A transaction verifier, in charge of verifying transaction occurrence by interfacing with the oracle. For more
// details, see transactionVerifier.go.
// This main function aims to demonstrate the interface between these two modules. In an actual
// light client, they can be entirely separate processes/smart contracts.
// The oracle should receive Algorand data from an off chain relayer, and the transaction verifier should receive
// transaction occurrence queries from third parties. For the purposes of this PoC, they're two separate go packages, and
// both the relayer and other third parties have been replaced with example committed data.
func main() {
	// This is the genesis data, required for initializing the oracle. This data can either be queried from the
	// blockchain itself, or found in the developer's portal.
	genesisVotersCommitment, genesisVotersLnProvenWeight, err := encoded_assets.GetParsedGenesisData("encoded_assets/genesis/")
	if err != nil {
		fmt.Printf("Failed to parse genesis assets: %s\n", err)
		return
	}

	// This is data required for verifying a transaction. In a real light client, this data should come from a
	// third party. The third party is responsible for querying Algorand to get most of this data.
	genesisHash, round, seed, transactionHash, transactionProofResponse, lightBlockHeaderProofResponse, err :=
		encoded_assets.GetParsedTransactionVerificationData("encoded_assets/transaction_verification/")
	if err != nil {
		fmt.Printf("Failed to parse assets needed for transaction verification: %s\n", err)
		return
	}

	// This is data required for advancing the oracle's state. In a real light client, this data should come from a relayer.
	stateProofMessage, stateProof, err :=
		encoded_assets.GetParsedStateProofAdvancmentData("encoded_assets/state_proof_verification/")
	if err != nil {
		fmt.Printf("Failed to parse assets needed for oracle state advancement: %s\n", err)
		return
	}

	// In a real light client, intervalSize and firstAttestedRound should be hardcoded, and retrieved from the Algorand
	// consensus and from the Algorand blockchain respectively.
	intervalSize := uint64(8)
	firstAttestedRound := uint64(9)

	// We initialize the oracle using the parsed genesis data and a hard coded capacity.
	oracleInstance := oracle.InitializeOracle(firstAttestedRound, intervalSize, genesisVotersCommitment, genesisVotersLnProvenWeight, 1000)

	// We advance the oracle's state using the state proof and the state proof message. The oracle verifies the message
	// using the state proof. See the documentation in oracle.go for more details.
	err = oracleInstance.AdvanceState(stateProof, stateProofMessage)
	if err != nil {
		fmt.Printf("Failed to advance oracle state: %s\n", err)
		return
	}

	// After advancing the oracle's state, we retrieve the commitment for the given transaction's round from the oracle.
	desiredTransactionCommitment, err := oracleInstance.GetStateProofCommitment(round)
	if err != nil {
		fmt.Printf("Failed to retrieve commitment interval for round %d: %s\n", round, err)
		return
	}

	// We then verify the transaction's occurrence using the data provided for transaction verification, along with the
	// commitment retrieved for the transaction.
	err = transactionverifier.VerifyTransaction(transactionHash, transactionProofResponse,
		lightBlockHeaderProofResponse, round, genesisHash, seed, desiredTransactionCommitment)

	if err != nil {
		fmt.Printf("Transaction verification failed: %s\n", err)
		return
	}
}
