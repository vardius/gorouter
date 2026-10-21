// Package benchdata holds the route set both routers are benchmarked against.
package benchdata

import "strings"

// Route is a pattern and a request path it matches.
type Route struct {
	Method  string
	Pattern string
	Path    string
}

// API is a REST API shaped like a large public one: mostly static and param routes, some regexps.
var API = build([]string{
	"GET /",
	"GET /emojis",
	"GET /events",
	"GET /feeds",
	"GET /meta",
	"GET /rate_limit",
	"GET /notifications",
	"GET /users",
	"GET /users/{user}",
	"GET /users/{user}/repos",
	"GET /users/{user}/followers",
	"GET /users/{user}/following",
	"GET /users/{user}/gists",
	"GET /users/{user}/orgs",
	"GET /users/{user}/events",
	"GET /orgs/{org}",
	"GET /orgs/{org}/repos",
	"POST /orgs/{org}/repos",
	"GET /orgs/{org}/members",
	"GET /orgs/{org}/members/{user}",
	"PUT /orgs/{org}/members/{user}",
	"DELETE /orgs/{org}/members/{user}",
	"GET /orgs/{org}/teams",
	"GET /repos/{owner}/{repo}",
	"PATCH /repos/{owner}/{repo}",
	"DELETE /repos/{owner}/{repo}",
	"GET /repos/{owner}/{repo}/issues",
	"POST /repos/{owner}/{repo}/issues",
	"GET /repos/{owner}/{repo}/issues/{number:[0-9]+}",
	"PATCH /repos/{owner}/{repo}/issues/{number:[0-9]+}",
	"GET /repos/{owner}/{repo}/issues/{number:[0-9]+}/comments",
	"POST /repos/{owner}/{repo}/issues/{number:[0-9]+}/comments",
	"GET /repos/{owner}/{repo}/pulls",
	"POST /repos/{owner}/{repo}/pulls",
	"GET /repos/{owner}/{repo}/pulls/{number:[0-9]+}",
	"GET /repos/{owner}/{repo}/pulls/{number:[0-9]+}/files",
	"PUT /repos/{owner}/{repo}/pulls/{number:[0-9]+}/merge",
	"GET /repos/{owner}/{repo}/commits",
	"GET /repos/{owner}/{repo}/commits/{sha}",
	"GET /repos/{owner}/{repo}/branches",
	"GET /repos/{owner}/{repo}/branches/{branch}",
	"GET /repos/{owner}/{repo}/releases",
	"GET /repos/{owner}/{repo}/releases/latest",
	"GET /repos/{owner}/{repo}/releases/{id:[0-9]+}",
	"GET /repos/{owner}/{repo}/contents/{path}",
	"GET /repos/{owner}/{repo}/stargazers",
	"GET /repos/{owner}/{repo}/forks",
	"GET /gists",
	"POST /gists",
	"GET /gists/public",
	"GET /gists/starred",
	"GET /gists/{id}",
	"PATCH /gists/{id}",
	"DELETE /gists/{id}",
	"GET /gists/{id}/comments",
	"GET /search/repositories",
	"GET /search/code",
	"GET /search/issues",
	"GET /search/users",
	"GET /user",
	"PATCH /user",
	"GET /user/repos",
	"POST /user/repos",
	"GET /user/emails",
	"GET /user/keys",
	"POST /user/keys",
	"GET /user/keys/{id}",
	"DELETE /user/keys/{id}",
	"PUT /user/starred/{owner}/{repo}",
	"DELETE /user/starred/{owner}/{repo}",
})

// build turns "METHOD /pattern" lines into routes, filling each param with a value that matches it.
func build(lines []string) []Route {
	routes := make([]Route, 0, len(lines))
	for _, line := range lines {
		method, pattern, _ := strings.Cut(line, " ")

		segments := strings.Split(pattern, "/")
		for i, s := range segments {
			switch {
			case strings.Contains(s, "[0-9]+"):
				segments[i] = "1347"
			case strings.HasPrefix(s, "{"):
				segments[i] = "gorouter"
			}
		}

		routes = append(routes, Route{Method: method, Pattern: pattern, Path: strings.Join(segments, "/")})
	}

	return routes
}
