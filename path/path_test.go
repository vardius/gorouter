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

func TestGetPart(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name         string
		args         args
		wantPart     string
		wantNextPath string
	}{
		{"empty", args{""}, "", ""},
		{"/", args{"/"}, "/", ""},
		{"x", args{"x"}, "x", ""},
		{"x/y", args{"x/y"}, "x", "y"},
		{"x/y/z", args{"x/y/z"}, "x", "y/z"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPart, gotNextPath := GetPart(tt.args.path)
			if gotPart != tt.wantPart {
				t.Errorf("GetPart() gotPart = %v, want %v", gotPart, tt.wantPart)
			}
			if gotNextPath != tt.wantNextPath {
				t.Errorf("GetPart() gotNextPath = %v, want %v", gotNextPath, tt.wantNextPath)
			}
		})
	}
}

func TestGetNameFromPart(t *testing.T) {
	type args struct {
		pathPart string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"x", args{"x"}, "x"},
		{"{name}", args{"{name}"}, "name"},
		{"{name:(w+)", args{"{name:(w+)"}, "name"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := GetNameFromPart(tt.args.pathPart)
			if got != tt.want {
				t.Errorf("GetNameFromPart() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStripLeadingSlashes(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		stripSlashes int
		want         string
	}{
		{"slashesCount = 0", "/foo/bar", 0, "/foo/bar"},
		{"slashesCount = 1", "/foo/bar", 1, "/bar"},
		{"slashesCount = 2", "/foo/bar", 2, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StripLeadingSlashes(tt.path, tt.stripSlashes); got != tt.want {
				t.Errorf("StripLeadingSlashes() = %v, want %v", got, tt.want)
			}
		})
	}
}
