package main

import (
	"crypto"
	"encoding/json"
	"fmt"
	"hash"
	"os"
)

func decodeFromFile(encodedPath string, target interface{}) error {
	encodedData, err := os.ReadFile(encodedPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(encodedData, target)
	return err
}

func hashBytes(hash hash.Hash, m []byte) []byte {
	hash.Reset()
	hash.Write(m)
	outhash := hash.Sum(nil)
	return outhash
}

func unmarshalHashFunc(hashStr string) (hash.Hash, error) {
	switch hashStr {
	case "sha256":
		return crypto.SHA256.New(), nil
	default:
		return nil, fmt.Errorf("unsupported hash function detected")
	}
}
