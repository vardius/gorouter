package goserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContext(t *testing.T) {
	req, err := http.NewRequest(get, "/x", nil)
	if err != nil {
		t.Fatal(err)
	}

	params := make(Params)
	params["test"] = "test"

			req = req.WithContext(newContextFromRequest(req, params))
	cParams := ParamsFromContext(req.Context())

	if params["test"] != cParams["test"] {
		t.Error("Request returned invalid context")
	}
}
