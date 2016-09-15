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
	case '[':
		a := []interface{}{}
		err := json.Unmarshal(b, &a)
		if err != nil {
			return nil, err
		}
		return a, nil
	default:
		return nil, errors.New("Unknown type in json bytes")
	}
}

func replaceMap(m map[string]interface{}, resolver func(string) (interface{}, error)) (interface{}, error) {
	if len(m) == 1 && m["Ref"] != nil {
		switch refval := m["Ref"].(type) {
		case string:
			return m, nil
		case map[string]interface{}:
			if len(refval) == 1 && refval["URI"] != nil {
				switch refurival := refval["URI"].(type) {
				case string:
					return resolver(refurival)
				default:
					return nil, errors.New("unexpected value for Ref:URI")
				}
			} else {
				replaced, err := replaceMap(refval, resolver)
				if err != nil {
					return nil, err
				}
				return map[string]interface{} { "Ref": replaced }, nil
			}
		default:
			return nil, errors.New("unknown reftype")
		}
	}
	returnable := make(map[string]interface{})
	for k, v := range m {
		val, err := replace(v, resolver)
		if err != nil {
			return nil, err
		}
		returnable[k] = val
	}
	return returnable, nil
}

func replace(i interface{}, resolver func(string) (interface{}, error)) (interface{}, error) {
	switch value := i.(type) {
	case string:
		return value, nil
	case map[string]interface{}:
		return replaceMap(value, resolver)
	default:
		return nil, nil
	}
}

func replaceBytes(b []byte, resolver func(string) (interface{}, error)) (interface{}, error) {
	i, err := bytesAsJson(b)
	if err != nil {
		return nil, err
	}

	return replace(i, resolver)
}

func main() {
}
