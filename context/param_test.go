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

func TestParamsSet(t *testing.T) {
	p := make(Params, 5)
	type args struct {
		index uint8
		key   string
		value string
	}
	tests := []struct {
		name string
		p    Params
		args args
	}{
		{"one", p, args{0, "one", "one"}},
		{"two", p, args{0, "two", "two"}},
		{"three", p, args{0, "three", "three"}},
		{"four", p, args{0, "four", "four"}},
		{"five", p, args{0, "five", "five"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Set(tt.args.index, tt.args.key, tt.args.value)

			if p.Value(tt.args.key) != tt.args.value {
				t.Error("Invalid parameter value")
			}
		})
	}
}
