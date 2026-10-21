package tree

import (
	"maps"
	"slices"
	"strings"
	"testing"

	"github.com/vardius/gorouter/v5/internal/race"
)

// lookup is a test helper returning the matched pattern and its params by name.
func lookup(t *testing.T, tr *Tree[string], method, path string) (string, map[string]string) {
	t.Helper()

	values := make([]string, tr.MaxParams())
	r := tr.Lookup(method, path, values)
	if r == nil {
		return "", nil
	}

	params := make(map[string]string, len(r.Params))
	for i, name := range r.Params {
		params[name] = values[i]
	}

	return r.Pattern, params
}

func TestLookup(t *testing.T) {
	t.Parallel()

	var tr Tree[string]
	for _, pattern := range []string{
		"/",
		"/users",
		"/users/new",
		"/users/{id:[0-9]+}",
		"/users/{name}",
		"/users/{name}/posts",
		"/users/{id:[0-9]+}/posts/{post}",
		"/a/b/c",
		"/a/{x}/d",
		"/colors/{c:red|green}",
		"/colors/{c:[a-z]+}",
	} {
		tr.Handle("GET", pattern, pattern)
	}
	tr.Mount("/v1", "mount")

	tests := []struct {
		path    string
		pattern string
		params  map[string]string
	}{
		{"/", "/", map[string]string{}},
		{"/users", "/users", map[string]string{}},
		{"/users/", "/users", map[string]string{}},
		{"/users/new", "/users/new", map[string]string{}},
		{"/users/42", "/users/{id:[0-9]+}", map[string]string{"id": "42"}},
		{"/users/bob", "/users/{name}", map[string]string{"name": "bob"}},
		{"/users/42/posts/7", "/users/{id:[0-9]+}/posts/{post}", map[string]string{"id": "42", "post": "7"}},
		// The regexp branch has no /posts route, so matching falls back to the wildcard.
		{"/users/42/posts", "/users/{name}/posts", map[string]string{"name": "42"}},
		{"/users/bob/posts", "/users/{name}/posts", map[string]string{"name": "bob"}},
		// Static b only leads to c, so /a/b/d falls back to {x}.
		{"/a/b/d", "/a/{x}/d", map[string]string{"x": "b"}},
		{"/a/b/c", "/a/b/c", map[string]string{}},
		// Regexps are tried in registration order and must match the whole segment.
		{"/colors/red", "/colors/{c:red|green}", map[string]string{"c": "red"}},
		{"/colors/reddish", "/colors/{c:[a-z]+}", map[string]string{"c": "reddish"}},
		{"/colors/RED", "", nil},
		{"/v1", "/v1", map[string]string{}},
		{"/v1/anything/below", "/v1", map[string]string{}},
		{"/v1x", "", nil},
		{"/users//posts", "", nil},
		// An empty segment never matches, not even right after the leading slash.
		{"//", "", nil},
		{"///", "", nil},
		{"/users//", "", nil},
		// OPTIONS * is about the server: it never reaches a route.
		{"*", "", nil},
		{"/nope", "", nil},
	}
	for _, tt := range tests {
		pattern, params := lookup(t, &tr, "GET", tt.path)
		if pattern != tt.pattern || !maps.Equal(params, tt.params) {
			t.Errorf("Lookup(GET, %q) = %q %v, want %q %v", tt.path, pattern, params, tt.pattern, tt.params)
		}
	}

	if pattern, _ := lookup(t, &tr, "POST", "/users"); pattern != "" {
		t.Errorf("Lookup(POST, /users) = %q, want no route", pattern)
	}
	for _, method := range []string{"POST", "PROPFIND"} {
		if pattern, _ := lookup(t, &tr, method, "/v1/x"); pattern != "/v1" {
			t.Errorf("Lookup(%s, /v1/x) = %q, want the mount on every method", method, pattern)
		}
	}
}

func TestMountCustomMethods(t *testing.T) {
	t.Parallel()

	var tr Tree[string]
	tr.Handle("REPORT", "/r", "report")
	tr.Mount("/dav", "dav")
	tr.Handle("PROPFIND", "/p", "propfind")
	tr.Use("MKCOL", "/dav", []func(string) string{func(h string) string { return "mkcol>" + h }})

	tests := []struct {
		method, path, want string
	}{
		{"REPORT", "/dav/x", "dav"},   // a custom method registered before the mount
		{"PROPFIND", "/dav/x", "dav"}, // one registered after it
		{"LOCK", "/dav/x", "dav"},     // one with no routes at all
		{"MKCOL", "/dav/x", "mkcol>dav"},
		{"LOCK", "/r", ""},
	}
	for _, tt := range tests {
		got := ""
		if r := tr.Lookup(tt.method, tt.path, nil); r != nil {
			got = r.Handler
		}
		if got != tt.want {
			t.Errorf("Lookup(%s, %s) = %q, want %q", tt.method, tt.path, got, tt.want)
		}
	}
}

func TestRoutesListsMountOnce(t *testing.T) {
	t.Parallel()

	var tr Tree[string]
	tr.Handle("GET", "/a", "")
	tr.Mount("/v1", "")
	tr.Handle("PROPFIND", "/b", "")

	var got []string
	for r := range tr.Routes() {
		got = append(got, r.Method+" "+r.Pattern)
	}
	if want := []string{"GET /a", " /v1", "PROPFIND /b"}; !slices.Equal(got, want) {
		t.Errorf("Routes() = %q, want %q", got, want)
	}
}

func TestValues(t *testing.T) {
	var tr Tree[string]
	tr.Handle("GET", "/{a}/{b}/{c}", "")

	v := tr.Values()
	if len(*v) != 3 {
		t.Fatalf("Values() holds %d values, want MaxParams 3", len(*v))
	}
	if r := tr.Lookup("GET", "/x/y/z", *v); r == nil || (*v)[2] != "z" {
		t.Fatalf("Lookup with Values() = %v %q", r, *v)
	}
	tr.PutValues(v)
	if race.Enabled {
		return // sync.Pool drops items at random under -race
	}
	if n := testing.AllocsPerRun(100, func() { tr.PutValues(tr.Values()) }); n != 0 {
		t.Errorf("Values and PutValues allocate %v times, want 0 once pooled", n)
	}
}

func TestLookupNamesParamsPerRoute(t *testing.T) {
	t.Parallel()

	// Both routes share the wildcard node; each names the value its own way.
	var tr Tree[string]
	tr.Handle("GET", "/users/{id}", "")
	tr.Handle("GET", "/users/{name}/posts", "")

	if _, params := lookup(t, &tr, "GET", "/users/7"); params["id"] != "7" {
		t.Errorf("params = %v, want id=7", params)
	}
	if _, params := lookup(t, &tr, "GET", "/users/bob/posts"); params["name"] != "bob" {
		t.Errorf("params = %v, want name=bob", params)
	}
}

func TestLookupCustomMethod(t *testing.T) {
	t.Parallel()

	var tr Tree[string]
	tr.Handle("PROPFIND", "/dav/{file}", "")

	if pattern, params := lookup(t, &tr, "PROPFIND", "/dav/a.txt"); pattern != "/dav/{file}" || params["file"] != "a.txt" {
		t.Errorf("Lookup(PROPFIND) = %q %v", pattern, params)
	}
}

func TestLookupDoesNotAllocate(t *testing.T) {
	var tr Tree[string]
	tr.Handle("GET", "/users/{id:[0-9]+}/posts/{post}", "")
	tr.Handle("GET", "/static/path/here", "")
	values := make([]string, tr.MaxParams())

	for _, path := range []string{"/users/42/posts/7", "/static/path/here", "/missing"} {
		if n := testing.AllocsPerRun(100, func() { tr.Lookup("GET", path, values) }); n != 0 {
			t.Errorf("Lookup(%q) allocates %v times, want 0", path, n)
		}
	}
}

func TestMiddleware(t *testing.T) {
	t.Parallel()

	tag := func(name string) func(string) string {
		return func(h string) string { return name + ">" + h }
	}

	var tr Tree[string]
	tr.Use("", "/", []func(string) string{tag("all")})
	tr.Handle("GET", "/x/y", "h")
	tr.Handle("POST", "/x/y", "h")
	tr.Handle("GET", "/xy", "h")
	tr.Handle("GET", "/users/{name}/posts", "h")
	tr.Handle("GET", "/users/{id:[0-9]+}", "h")
	tr.Mount("/v1", "sub")
	tr.Use("GET", "/x", []func(string) string{tag("get-x")})
	tr.Use("", "/x/y", []func(string) string{tag("first"), tag("second")})
	tr.Use("", "/users/{id}", []func(string) string{tag("wildcard")})
	tr.Use("", "/v1", []func(string) string{tag("v1")})
	tr.Use("", "/z", []func(string) string{tag("z")})

	tests := []struct {
		method, path, want string
	}{
		{"GET", "/x/y", "all>get-x>first>second>h"},
		{"POST", "/x/y", "all>first>second>h"},
		{"GET", "/xy", "all>h"},
		// Param names don't matter, the node does: {id} and {name} are the same wildcard.
		{"GET", "/users/bob/posts", "all>wildcard>h"},
		// A regexp param is a different node from the wildcard.
		{"GET", "/users/7", "all>h"},
		{"DELETE", "/v1/a/b", "all>v1>sub"},
	}
	for _, tt := range tests {
		r := tr.Lookup(tt.method, tt.path, make([]string, tr.MaxParams()))
		if r == nil || r.Handler != tt.want {
			t.Errorf("%s %s handler = %v, want %q", tt.method, tt.path, r, tt.want)
		}
	}
}

func TestAllowed(t *testing.T) {
	t.Parallel()

	var tr Tree[string]
	tr.Handle("DELETE", "/x", "")
	tr.Handle("POST", "/x", "")
	tr.Handle("GET", "/x", "")
	tr.Handle("PROPFIND", "/x", "")
	tr.Handle("OPTIONS", "/x", "")
	tr.Handle("PUT", "/", "")

	tests := []struct {
		method, path, want string
	}{
		{"PATCH", "/x", "GET, POST, DELETE, PROPFIND, OPTIONS"},
		{"GET", "/x/", "POST, DELETE, PROPFIND, OPTIONS"},
		{"GET", "/", "PUT, OPTIONS"},
		{"OPTIONS", "*", "GET, POST, PUT, DELETE, PROPFIND, OPTIONS"},
		{"GET", "/nope", ""},
	}

	// UseMethod gives POST a tree, but with no POST routes the server doesn't allow it.
	var unused Tree[string]
	unused.Handle("GET", "/a", "")
	unused.Use("POST", "/a", []func(string) string{func(h string) string { return h }})
	if got := unused.Allowed("OPTIONS", "*", nil); got != "GET, OPTIONS" {
		t.Errorf("Allowed(OPTIONS, *) with POST middleware but no POST routes = %q, want GET, OPTIONS", got)
	}

	// A mount whose first segment is a param would match "*" as a path; OPTIONS * must not reach it.
	var mounted Tree[string]
	mounted.Handle("POST", "/login", "")
	mounted.Mount("/{tenant}", "")
	if r := mounted.Lookup("OPTIONS", "*", make([]string, 1)); r != nil {
		t.Errorf("Lookup(OPTIONS, *) = %s, want no route", r.Pattern)
	}
	for _, tt := range tests {
		if got := tr.Allowed(tt.method, tt.path, make([]string, tr.MaxParams())); got != tt.want {
			t.Errorf("Allowed(%s, %q) = %q, want %q", tt.method, tt.path, got, tt.want)
		}
	}
}

func TestRegistrationPanics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		register func(*Tree[string])
		want     string
	}{
		{"no leading slash", func(tr *Tree[string]) { tr.Handle("GET", "users", "") }, "must start with /"},
		{"empty segment", func(tr *Tree[string]) { tr.Handle("GET", "/a//b", "") }, "empty segment"},
		{"empty param name", func(tr *Tree[string]) { tr.Handle("GET", "/{}", "") }, "needs a name"},
		{"unclosed param", func(tr *Tree[string]) { tr.Handle("GET", "/{id", "") }, "not closed"},
		{"regexp with a slash", func(tr *Tree[string]) { tr.Handle("GET", "/{p:a/b}", "") }, "not closed"},
		{"text and param", func(tr *Tree[string]) { tr.Handle("GET", "/a{id}", "") }, "mixes text"},
		{"empty regexp", func(tr *Tree[string]) { tr.Handle("GET", "/{id:}", "") }, "empty regexp"},
		{"bad regexp", func(tr *Tree[string]) { tr.Handle("GET", "/{id:[}", "") }, "missing closing ]"},
		{"repeated param", func(tr *Tree[string]) { tr.Handle("GET", "/{a}/{a}", "") }, "repeats param {a}"},
		{"empty method", func(tr *Tree[string]) { tr.Handle("", "/", "") }, "empty method"},
		{"duplicate", func(tr *Tree[string]) {
			tr.Handle("GET", "/users/{id}", "")
			tr.Handle("GET", "/users/{name}/", "")
		}, "GET /users/{name}/ conflicts with GET /users/{id}"},
		{"route under a mount", func(tr *Tree[string]) {
			tr.Mount("/v1", "")
			tr.Handle("GET", "/v1/users", "")
		}, "conflicts with the mount at /v1"},
		{"mount over routes", func(tr *Tree[string]) {
			tr.Handle("GET", "/v1/users", "")
			tr.Mount("/v1", "")
		}, "mount at /v1 conflicts with GET routes"},
		{"middleware below a mount", func(tr *Tree[string]) {
			tr.Mount("/api/v1", "")
			tr.Use("", "/api/v1/admin", []func(string) string{func(h string) string { return h }})
		}, "middleware at /api/v1/admin is below the mount at /api/v1"},
		{"mount over middleware", func(tr *Tree[string]) {
			tr.Use("GET", "/{tenant}/admin", []func(string) string{func(h string) string { return h }})
			tr.Mount("/{t}", "")
		}, "mount at /{t} hides the middleware used below it"},
		{"custom route under a mount", func(tr *Tree[string]) {
			tr.Mount("/dav", "")
			tr.Handle("PROPFIND", "/dav/x", "")
		}, "PROPFIND /dav/x conflicts with the mount at /dav"},
		{"bad middleware pattern", func(tr *Tree[string]) {
			tr.Use("", "x", []func(string) string{func(h string) string { return h }})
		}, "must start with /"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				msg, _ := recover().(string)
				if !strings.Contains(msg, tt.want) {
					t.Errorf("panic = %q, want it to contain %q", msg, tt.want)
				}
			}()
			tt.register(new(Tree[string]))
		})
	}
}

func TestStripSegments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		path string
		n    int
		want string
	}{
		{"/v1/users/7", 1, "/users/7"},
		{"/v1/users/7", 2, "/7"},
		{"/v1/users/7", 3, "/"},
		{"/v1", 1, "/"},
		{"/v1/", 1, "/"},
		{"/v1/x", 0, "/v1/x"},
		{"", 0, "/"},
	}
	for _, tt := range tests {
		if got := StripSegments(tt.path, tt.n); got != tt.want {
			t.Errorf("StripSegments(%q, %d) = %q, want %q", tt.path, tt.n, got, tt.want)
		}
	}
}

func TestDepth(t *testing.T) {
	t.Parallel()

	for pattern, want := range map[string]int{"/": 0, "/v1": 1, "/v1/": 1, "/a/{b}/c": 3} {
		if got := Depth(pattern); got != want {
			t.Errorf("Depth(%q) = %d, want %d", pattern, got, want)
		}
	}
}

func TestParamNames(t *testing.T) {
	t.Parallel()

	if got := ParamNames("/a/{b}/c/{d:[0-9]+}"); !slices.Equal(got, []string{"b", "d"}) {
		t.Errorf("ParamNames = %q, want [b d]", got)
	}
}

func TestEmptyTree(t *testing.T) {
	t.Parallel()

	var tr Tree[string]
	tr.Use("", "/", nil) // no middleware: nothing to register
	if r := tr.Lookup("GET", "/", nil); r != nil {
		t.Errorf("Lookup on an empty tree = %s, want nil", r.Pattern)
	}

	tr.Handle("GET", "/a", "")
	tr.Handle("GET", "/b", "")
	n := 0
	for range tr.Routes() {
		n++
		break
	}
	if n != 1 {
		t.Errorf("Routes yielded %d routes after break, want 1", n)
	}
}
