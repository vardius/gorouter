package goserver

import "strings"

type paths []string

func getPaths(path string) paths {
	if path != "" && path[0] == '/' {
		path = path[1:]
	}
	pl := len(path)
	if path != "" && path[pl-1] == '/' {
		path = path[:pl-1]
	}

	var p paths
	pathsLen := 0
	lastIndex := 0
	if len(path) > 0 {
		for i := 0; i < len(path); i++ {
			if path[i] == '{' {
				if i == len(path)-1 && lastIndex == 0 {
					p.add(strings.Trim(path[:], "/"))
				} else if lastIndex == 0 {
					p.add(strings.Trim(path[:i], "/"))
				} else if i == len(path)-1 {
					p.add(strings.Trim(path[lastIndex-1:], "/"))
				} else {
					p.add(strings.Trim(path[lastIndex-1:i], "/"))
				}
				pathsLen++
				lastIndex = i
			} else if path[i] == '}' {
				if i == len(path)-1 && lastIndex == 0 {
					p.add(strings.Trim(path[:], "/"))
				} else if lastIndex == 0 {
					p.add(strings.Trim(path[:i+1], "/"))
				} else if i == len(path)-1 {
					p.add(strings.Trim(path[lastIndex:], "/"))
				} else {
					p.add(strings.Trim(path[lastIndex:i+1], "/"))
				}
				pathsLen++
				lastIndex = i
			}
		}
		if path[len(path)-1] != '}' {
			p.add(strings.Trim(path[lastIndex+1:], "/"))
		}
	}
	return p
}

func (p paths) add(part string) paths {
	if part != "" && part != "}" && part != "{" && part != "{}" {
		p = append(p, part)
	}
	return p
}
