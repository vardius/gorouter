package gorouter

import (
	"net/http"
	"testing"
)

func TestContext(t *testing.T) {
	req, err := http.NewRequest(GET, "/x", nil)
	if err != nil {
		t.Fatal(err)
	}

	p := Param{"test", "test"}
	params := Params{p}

	req = req.WithContext(newContext(req, params))
	cParams, ok := FromContext(req.Context())
	if !ok {
		t.Fatal("Error while getting context")
	}

	if params.Value("test") != cParams.Value("test") {
		t.Error("Request returned invalid context")
	}
}
