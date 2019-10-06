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

func GetNameFromPart(pathPart string) (id string, exp string) {
	id = pathPart

	if pathPart[0] == '{' {
		id = pathPart[1 : len(pathPart)-1]

		if parts := strings.Split(id, ":"); len(parts) == 2 {
			id = parts[0]
			exp = parts[1]
		}

		if id == "" {
			panic("Empty wildcard name")
		}

		return
	}

	return
}
