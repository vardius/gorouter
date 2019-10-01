package path

import (
	"testing"
)

func TestTrim(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{""}, ""},
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
			if got := TrimSlash(tt.args.path); got != tt.want {
				t.Errorf("[%s] Trim() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
