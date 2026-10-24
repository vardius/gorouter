module github.com/vardius/gorouter/examples/grpc

go 1.26.0

require (
	connectrpc.com/connect v1.21.1
	connectrpc.com/grpchealth v1.5.0
	github.com/vardius/gorouter/v5 v5.0.0
	google.golang.org/grpc v1.84.0
)

require (
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260706201446-f0a921348800 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/vardius/gorouter/v5 => ../..
