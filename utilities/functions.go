package utilities

import (
	"crypto"
	"encoding/json"
	"fmt"
	"hash"
	"os"
)

func DecodeFromFile(encodedPath string, target interface{}) error {
	encodedData, err := os.ReadFile(encodedPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(encodedData, target)
	return err
}

func HashBytes(hash hash.Hash, m []byte) []byte {
	hash.Reset()
	hash.Write(m)
	outhash := hash.Sum(nil)
	return outhash
}

func UnmarshalHashFunc(hashStr string) (hash.Hash, error) {
	switch hashStr {
	case "sha256":
		return crypto.SHA256.New(), nil
	default:
		return nil, fmt.Errorf("unsupported hash function detected")
	}
}
