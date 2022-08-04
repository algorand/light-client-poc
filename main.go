package main

import (
	"fmt"
	"github.com/almog-t/light-client-poc/encoded_assets"
	"github.com/almog-t/light-client-poc/oracle"
	"github.com/almog-t/light-client-poc/transactionverifier"
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
	firstAttestedRound := stateProofMessage.FirstAttestedRound
	oracleInstance := oracle.InitializeOracle(firstAttestedRound, intervalSize, genesisVotersCommitment, genesisVotersLnProvenWeight, 1000)
	err = oracleInstance.AdvanceState(stateProof, stateProofMessage)
	if err != nil {
		fmt.Printf("Failed to advance oracle state: %s\n", err)
		return
	}

	desiredTransactionCommitment, err := oracleInstance.GetStateProofCommitment(round)
	if err != nil {
		fmt.Printf("Failed to retrieve commitment interval for round %d: %s\n", round, err)
		return
	}

	err = transactionverifier.VerifyTransaction(transactionHash, transactionProofResponse,
		lightBlockHeaderProofResponse, round, genesisHash, seed, desiredTransactionCommitment)

	if err != nil {
		fmt.Printf("Transaction verification failed: %s\n", err)
		return
	}
}
