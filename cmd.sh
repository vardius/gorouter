#!/bin/sh

go get -u -v github.com/derekparker/delve/cmd/dlv
go get -u -v github.com/stretchr/testify/assert

go fix
go fmt
go build
go vet
go test
dlv test --headless --listen=:2345 --log
