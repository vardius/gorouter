#!/bin/sh

go get -u github.com/golang/dep/...
# dep init
dep ensure -update

# dep ensure github.com/pkg/errors@^0.8.0

# go get -u -v github.com/derekparker/delve/cmd/dlv
go get -u -v github.com/stretchr/testify/assert

go fix
go fmt
go build
go vet
go test
# dlv test --headless --listen=:2345 --log
