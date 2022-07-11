package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

func getTransactionProofResponse() (models.ProofResponse, error) {
	proofJsonStr := `{
"hashtype": "sha256",
"idx": 0,
"proof": "5xvqtS3zHjSUWiBUK8YupN8AXvl733os1Xln8RfEZBA=",
"stibhash": "f8mCwgwSThJwH37N5hBj6CjdO/ZhH8XXby2/Bke1lN0=",
"treedepth": 1
}`
	jsonBytes := []byte(proofJsonStr)
	resp := models.ProofResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	return resp, err
}

func getBlockIntervalCommitment() ([]byte, error) {
	commitmentStr := "0QgCvCujCapNmpxTiVv5meq3WLfA3R8G857/nF0iyF0="
	return base64.StdEncoding.DecodeString(commitmentStr)
}

func getTransactionId() ([]byte, error) {
	transactionIdEncoded := "DQA86GJmEfXKpOCFtbW31EYoFqSjRR8/t63RGCajkHA="
	return base64.StdEncoding.DecodeString(transactionIdEncoded)
}

func getLightBlockHeaderProof() (LightBlockHeaderProof, error) {
	proofJsonStr := `{
"index": 0,
"proof": "oNr1sknaf3Hb9rohvNMVRt+LVg71Q3bQa+Fn7u9IRDoQEwCtQSIdOKhrmmubtJq5l7PFaH452Og/xPKqCKySCje8rezV8J56Znxge8MTaF66c6NYOKbgrDq7OvCUiYjX",
"treedepth": 3
}`
	jsonBytes := []byte(proofJsonStr)
	resp := LightBlockHeaderProof{}

	err := json.Unmarshal(jsonBytes, &resp)
	return resp, err
}

func getRound() uint64 {
	return 9
}

func getGenesisHash() ([]byte, error) {
	genesisHashEncoded := "SvwlI5kslp0rKzWBFTxeIdp6nxxxZa97LJx03F39bEQ="
	return base64.StdEncoding.DecodeString(genesisHashEncoded)
}

func TestTransactionVerify(t *testing.T) {
	r := require.New(t)
	blockIntervalCommitment, err := getBlockIntervalCommitment()
	r.NoError(err)

	transactionId, err := getTransactionId()
	r.NoError(err)

	transactionProof, err := getTransactionProofResponse()
	r.NoError(err)

	lightBlockHeaderProof, err := getLightBlockHeaderProof()
	r.NoError(err)

	genesisHash, err := getGenesisHash()
	r.NoError(err)

	round := getRound()

	verified, err := VerifyTransaction(transactionId, transactionProof,
		lightBlockHeaderProof, blockIntervalCommitment, genesisHash, round)
	r.NoError(err)
	if !verified {
		t.Fail()
	}
}
