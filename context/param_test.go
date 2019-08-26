package context

import (
	"testing"
)

func TestParamValue(t *testing.T) {
	key := "key"
	value := "value"

	p := Param{key, value}
	params := Params{p}

	if params.Value(key) != value {
		t.Error("Invalid parameter value")
	}
}
