package goserver

import (
	"net/http"
	"testing"
)

func TestContext(t *testing.T) {
	req, err := http.NewRequest(GET, "/x", nil)
	if err != nil {
		t.Fatal(err)
	}

	params := make(Params)
	params["test"] = "test"

	req = req.WithContext(newContextFromRequest(req, params))
	cParams, ok := ParamsFromContext(req.Context())
	if !ok {
		t.Fatal("Error while getting context")
	}

	if params["test"] != cParams["test"] {
		t.Error("Request returned invalid context")
	}
}
