package main

import (
	"encoding/json"
	"testing"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

func TestTransactionVerify(t *testing.T) {
	json_str := `{
"hashtype": "sha256",
"idx": 5,
"proof": "aiX/NAh7vN4w/jUeBNK2IIL2CbeSZ+rh+/J5fIp4SiabNrChkdR2DDP4QiS5FL5LXkkZ0H3vpODwDIDl0YZ33pgv8bjmkF3WWxl8T+ra2ab9zteBHT8MGWd8qheCPmtC",
"stibhash": "uxH+j9TV9MtgnKCQ4RGBZBscKjFry/bOOYF1VBdtBAg=",
"treedepth": 3
}`
	jsonBytes := []byte(json_str)
	resp := models.ProofResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		t.Fail()
	}

	commitmentStr := `[
7,
144,
104,
105,
180,
90,
77,
153,
137,
213,
221,
218,
90,
2,
125,
155,
146,
170,
101,
250,
198,
172,
96,
152,
95,
195,
64,
204,
147,
102,
23,
88
]`
	commitmentStrBytes := []byte(commitmentStr)
	commitment := [32]byte{}
	err = json.Unmarshal(commitmentStrBytes, &commitment)
	if err != nil {
		t.Fail()
	}

	idStr := `[
110,
107,
58,
8,
150,
136,
48,
225,
11,
24,
220,
98,
72,
132,
148,
229,
241,
232,
191,
202,
218,
11,
143,
79,
207,
146,
178,
122,
142,
162,
37,
183
]`

	txIdBytes := []byte(idStr)
	txId := [32]byte{}
	err = json.Unmarshal(txIdBytes, &txId)
	if err != nil {
		t.Fail()
	}

	verified, err := verifyTransaction(commitment[:], txId[:], resp)
	if !verified || err != nil {
		t.Fail()
	}
}
