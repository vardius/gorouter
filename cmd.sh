#!/bin/sh

go fix
go fmt
go build
go vet
go test
go test -bench=. -cpu=4 
