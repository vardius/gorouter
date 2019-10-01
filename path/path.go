package path

import "strings"

// Split splits URL path into parts
func Split(path string) (parts []string) {
	for {
		if i := strings.IndexByte(path, '{'); i >= 0 {
			if part := TrimSlash(path[:i]); part != "" {
				parts = append(parts, part)
			}
			if j := strings.IndexByte(path, '}') + 1; j > 0 {
				if part := TrimSlash(path[i:j]); part != "" {
					parts = append(parts, part)
				}
				i = j
			} else {
				continue
			}
			path = path[i:]
		} else {
			break
		}
	}

	if len(path) != 0 && path != "/" {
		if part := TrimSlash(path); part != "" {
			parts = append(parts, part)
		}
	}

	return
}

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
