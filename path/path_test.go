package path

import (
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name      string
		args      args
		wantParts []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotParts := Split(tt.args.path); !reflect.DeepEqual(gotParts, tt.wantParts) {
				t.Errorf("Split() = %v, want %v", gotParts, tt.wantParts)
			}
		})
	}
}

func TestTrim(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Trim(tt.args.path); got != tt.want {
				t.Errorf("Trim() = %v, want %v", got, tt.want)
			}
		})
	}
}
