FROM golang:latest

RUN go get -u -v github.com/derekparker/delve/cmd/dlv

WORKDIR /go/src/github.com/vardius/goserver

RUN go fix
RUN go fmt
RUN go build
RUN go vet
RUN go test
