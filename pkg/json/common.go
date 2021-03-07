package json

import (
	"encoding/json"
	"net/http"
)

// Decode decode json byte array to type
func Decode(v interface{}, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
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
