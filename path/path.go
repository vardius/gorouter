package path

import "strings"

// TrimSlash trims '/' URL path
func TrimSlash(path string) string {
	pathLen := len(path)
	if pathLen > 0 && path[0] == '/' {
		path = path[1:]
		pathLen--
	}

	if pathLen > 0 && path[pathLen-1] == '/' {
		path = path[:pathLen-1]
	}

	return path
}

// GetPart returns first path part and next path as a second argument
func GetPart(path string) (part string, nextPath string) {
	if j := strings.IndexByte(path, '/'); j > 0 {
		part = path[:j]
		nextPath = path[j+1:]
	} else {
		part = path
		nextPath = ""
	}

	return
}

// GetNameFromPart gets node name from path part, braces are not stripped so that we can check which node is associated with a parameter
func GetNameFromPart(pathPart string) (name string, exp string) {
	name = pathPart

	if parts := strings.Split(name, ":"); len(parts) == 2 {
		name = parts[0]
		exp = parts[1]
	}

	if name == "" {
		panic("Empty wildcard name")
	}

	return
}
