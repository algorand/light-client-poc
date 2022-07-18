package utilities

import (
	"encoding/json"
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
