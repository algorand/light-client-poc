package main

import (
	"fmt"
	"github.com/almog-t/light-client-poc/encoded_assets"
	"github.com/almog-t/light-client-poc/light_client_components"
)

func main() {
	genesisVotersCommitment, genesisVotersLnProvenWeight, err := encoded_assets.GetParsedGenesisData("encoded_assets/genesis/")
	if err != nil {
		fmt.Printf("Failed to parse genesis assets: %s\n", err)
		return
	}

	genesisHash, round, seed, transactionHash, transactionProofResponse, lightBlockHeaderProofResponse, err :=
		encoded_assets.GetParsedTransactionVerificationData("encoded_assets/transaction_verification/")
	if err != nil {
		fmt.Printf("Failed to parse assets needed for transaction verification: %s\n", err)
		return
	}

	stateProofMessage, stateProof, err :=
		encoded_assets.GetParsedStateProofAdvancmentData("encoded_assets/state_proof_verification/")
	if err != nil {
		fmt.Printf("Failed to parse assets needed for oracle state advancement: %s\n", err)
		return
	}

	intervalSize := stateProofMessage.LastAttestedRound - stateProofMessage.FirstAttestedRound + 1
	oracle := light_client_components.InitializeOracle(intervalSize, genesisVotersCommitment, genesisVotersLnProvenWeight, 1000)
	err = oracle.AdvanceState(stateProof, stateProofMessage)
	if err != nil {
		fmt.Printf("Failed to advance oracle state: %s\n", err)
		return
	}

	desiredTransactionCommitment, err := oracle.GetStateProofCommitment(round)
	if err != nil {
		fmt.Printf("Failed to retrieve commitment interval for round %d: %s\n", round, err)
		return
	}

	err = light_client_components.VerifyTransaction(transactionHash, transactionProofResponse,
		lightBlockHeaderProofResponse, round, genesisHash, seed, desiredTransactionCommitment)

	if err != nil {
		fmt.Printf("Transaction verification failed: %s\n", err)
		return
	}
}
