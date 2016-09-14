package main

import (
	"bytes"
	"encoding/json"
	"errors"
)

func bytesAsJson(b []byte) (interface{}, error) {
	b = bytes.TrimSpace(b)
	switch b[0] {
	case '"':
		s := ""
		err := json.Unmarshal(b, &s)
		if err != nil {
			return nil, err
		}
		return s, nil
	case '{':
		m := map[string]interface{}{}
		err := json.Unmarshal(b, &m)
		if err != nil {
			return nil, err
		}
		return m, nil
	default:
		return nil, errors.New("Unknown type")
	}
}

func main() {
}
