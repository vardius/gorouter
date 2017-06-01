#!/bin/sh

go get github.com/julienschmidt/httprouter

go fix
go fmt
go build
go vet
go test

# go test -run benchmark_gorouter_test.go -bench="BenchmarkGoRouter*" -cpu=4
# go test -run benchmark_httprouter_test.go -bench="BenchmarkHttpRouter*" -cpu=4
