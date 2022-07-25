package main

import (
	"fmt"
	"light-client-poc/encoded_assets"
)

func main() {
	genesis_hash, round, seed, transactionHash, transactionProofResponse, lightBlockHeaderProofResponse, _, err := encoded_assets.GetParsedTransactionVerificationData()
	if err != nil {
		fmt.Printf("Failed to parse assets needed for transaction verification: %s\n", err)
		return
	}

	genesisVotersCommitment, genesisVotersLnProvenWeight, stateProofMessage, stateProof, err := encoded_assets.GetParsedStateProofAdvancmentData()
	if err != nil {
		fmt.Printf("Failed to parse assets needed for oracle state advancement: %s\n", err)
		return
	}

	intervalSize := stateProofMessage.LastAttestedRound - stateProofMessage.FirstAttestedRound + 1
	// Actual initialization will not use stateProofMessage.FirstAttestedRound, it's simply convenient here for testing
	// purposes.
	oracle := InitializeOracle(intervalSize, stateProofMessage.FirstAttestedRound, genesisVotersCommitment, genesisVotersLnProvenWeight)
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

	transactionVerifier := TransactionVerifier{genesisHash: genesis_hash}
	err = transactionVerifier.VerifyTransaction(transactionHash, transactionProofResponse,
		lightBlockHeaderProofResponse, round, seed, desiredTransactionCommitment)

	if err != nil {
		fmt.Printf("Transaction verification failed: %s\n", err)
		return
	}
}
