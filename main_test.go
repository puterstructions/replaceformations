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

func equivalentString(value string, expected interface{}, t *testing.T) {
	switch v2 := expected.(type) {
	case string:
		if value != v2 {
			t.Errorf("%s is not an equivalent string to %s", value, expected)
		}
	default:
		t.Errorf("%v is not a string", expected)
	}
}

func equivalentMap(value map[string]interface{}, expected interface{}, t *testing.T) {
	switch v2 := expected.(type) {
	case map[string]interface{}:
		if len(value) != len(v2) {
			t.Errorf("len(%v) != len(%v)", value, v2)
		}
		for k, v := range value {
			equivalentJson(v, v2[k], t)
		}
	default:
		t.Errorf("%v is not a map", expected)
	}
}

func equivalentJson(value interface{}, expected interface{}, t *testing.T) {
	switch v := value.(type) {
	case string:
		equivalentString(v, expected, t)
	case map[string]interface{}:
		equivalentMap(v, expected, t)
	default:
		t.Errorf("%v is not equivalent to %v", value, expected)
	}
}

func asJson(raw []byte, t *testing.T) interface{} {
	v, err := bytesAsJson(raw)
	if err != nil {
		t.Error(err)
	}
	return v
}

func TestEquivalent(t *testing.T) {
	equivalentJson(
		asJson([]byte(`"asdf"`), t),
		asJson([]byte(`"asdf"`), t),
		t,
	)
	equivalentJson(
		asJson([]byte(`{"foo":"bar"}`), t),
		asJson([]byte(`{"foo":"bar"}`), t),
		t,
	)
}

func resolver(t *testing.T) func(string) (interface{}, error) {
	return func(name string) (interface{}, error) {
		return asJson([]byte(`{"a":"thing"}`), t), nil
	}
}

func TestIdentityReplaceString(t *testing.T) {
	replaced, err := replace(asJson([]byte(`"zxcvb"`), t), resolver(t))
	if err != nil {
		t.Error(err)
	}

	equivalentJson(
		replaced,
		asJson([]byte(`"zxcvb"`), t),
		t,
	)
}

func TestIdentityReplaceAMap(t *testing.T) {
	template := `{"asdf":"qwerty"}`
	replaced, err := replace(asJson([]byte(template), t), resolver(t))
	if err != nil {
		t.Error(err)
	}

	expected := asJson([]byte(template), t)
	equivalentJson(
		replaced,
		expected,
		t,
	)
}

func TestReplaceMapIdentity(t *testing.T) {
	template := map[string]interface{} { "asdf": "qwerty" }
	replaced, err := replaceMap(template, resolver(t))
	if err != nil {
		t.Error(err)
	}

	expected := asJson([]byte(`{"asdf":"qwerty"}`), t)
	equivalentJson(
		replaced,
		expected,
		t,
	)
}

func TestReplaceMapRefIdentity(t *testing.T) {
	template := map[string]interface{} { "Ref": "foo" }
	replaced, err := replaceMap(template, resolver(t))
	if err != nil {
		t.Error(err)
	}

	expected := asJson([]byte(`{"Ref":"foo"}`), t)
	equivalentJson(
		replaced,
		expected,
		t,
	)
}

func TestReplaceMapRefURI(t *testing.T) {
	template := map[string]interface{} { "Ref": map[string]interface{} { "URI": "somecomponent" } }
	replaced, err := replaceMap(template, resolver(t))
	if err != nil {
		t.Error(err)
	}

	expected := asJson([]byte(`{"a":"thing"}`), t)
	equivalentJson(
		replaced,
		expected,
		t,
	)
}

func TestIdentityRef(t *testing.T) {
	template := `{"asdf":{"Ref": "NamedValue"}}`
	replaced, err := replace(asJson([]byte(template), t), resolver(t))
	if err != nil {
		t.Error(err)
	}

	expected := asJson([]byte(template), t)
	equivalentJson(
		replaced,
		expected,
		t,
	)
}

func TestReplace(t *testing.T) {
	template := `{"asdf":{"Ref": {"URI":"componentName"}}}`
	replaced, err := replace(asJson([]byte(template), t), resolver(t))
	if err != nil {
		t.Error(err)
	}

	expected := asJson([]byte(`{"asdf":{"a":"thing"}}`), t)
	equivalentJson(
		replaced,
		expected,
		t,
	)
}
