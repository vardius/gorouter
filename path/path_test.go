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
		{"/", args{"/"}, nil},
		{"/x", args{"/x"}, []string{"x"}},
		{"/x/", args{"/x/"}, []string{"x"}},
		{"/x/y", args{"/x/y"}, []string{"x/y"}},
		{"/x/y/", args{"/x/y/"}, []string{"x/y"}},
		{"/{name}", args{"/{name}"}, []string{"{name}"}},
		{"/{name}/", args{"/{name}/"}, []string{"{name}"}},
		{"/x/{name}", args{"/x/{name}"}, []string{"x", "{name}"}},
		{"/x/{name}/", args{"/x/{name}/"}, []string{"x", "{name}"}},
		{"/x/{name}/y", args{"/x/{name}/y"}, []string{"x", "{name}", "y"}},
		{"/x/{name}/y/", args{"/x/{name}/y/"}, []string{"x", "{name}", "y"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotParts := Split(tt.args.path); !reflect.DeepEqual(gotParts, tt.wantParts) {
				t.Errorf("[%s] Split() = %v, want %v", tt.name, gotParts, tt.wantParts)
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
		{"/", args{"/"}, ""},
		{"/x", args{"/x"}, "x"},
		{"/x/", args{"/x/"}, "x"},
		{"/x/y", args{"/x/y"}, "x/y"},
		{"/x/y/", args{"/x/y/"}, "x/y"},
		{"/{name}", args{"/{name}"}, "{name}"},
		{"/{name}/", args{"/{name}/"}, "{name}"},
		{"/x/{name}", args{"/x/{name}"}, "x/{name}"},
		{"/x/{name}/", args{"/x/{name}/"}, "x/{name}"},
		{"/x/{name}/y", args{"/x/{name}/y"}, "x/{name}/y"},
		{"/x/{name}/y/", args{"/x/{name}/y/"}, "x/{name}/y"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Trim(tt.args.path); got != tt.want {
				t.Errorf("[%s] Trim() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
