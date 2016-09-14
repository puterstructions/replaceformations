package main

import (
	"testing"
)

func TestBytesAsJsonString(t *testing.T) {
	value, err := bytesAsJson([]byte(`"qwerty"`))
	if err != nil {
		t.Error(err)
	}

	if value.(string) != "qwerty" {
		t.Errorf("'%v' is not querty", value)
	}
}

func TestBytesAsJsonMap(t *testing.T) {
	value, err := bytesAsJson([]byte(`{"qwerty":"uiopas"}`))
	if err != nil {
		t.Error(err)
	}

	m := value.(map[string]interface{})
	if len(m) != 1 {
		t.Error("not length 1")
	}

	v := m["qwerty"]
	if v.(string) != "uiopas" {
		t.Errorf("'%v' is not uiopas", v)
	}
}
