package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"cudi/internal/pkg/domain"
)

func ParseConfig(filePath string) (domain.Cleanup, error) {
	var cleanup domain.Cleanup
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return cleanup, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return cleanup, err
	}
	return cleanup, json.Unmarshal(byteValue, &cleanup)
}

