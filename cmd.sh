#!/bin/sh

go fix
go fmt
go build
go vet
go test
