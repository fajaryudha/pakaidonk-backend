package utils

import (
	"bytes"
	"encoding/json"
	"log"
)

// MinifyJSON removes unnecessary spaces and newlines from a JSON object
func MinifyJSON(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := json.Compact(&buf, input)
	if err != nil {
		log.Printf("Error minifying JSON: %v", err)
		return nil, err
	}
	return buf.Bytes(), nil
}
