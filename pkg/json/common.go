package json

import (
	"encoding/json"
	"io"
)

// Decode decode json byte array to type
func Decode(v interface{}, readerCloser io.ReadCloser) error {
	decoder := json.NewDecoder(readerCloser)
	return decoder.Decode(v)
}

// Marshal marshall JSON to byte array
func Marshal(v interface{}) ([]byte, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func Unmarshal(source string, v interface{}) error {
	return json.Unmarshal([]byte(source), v)
}
