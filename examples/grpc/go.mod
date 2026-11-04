module github.com/vardius/gorouter/examples/grpc

go 1.26.0

require (
	connectrpc.com/connect v1.21.1
	connectrpc.com/grpchealth v1.5.0
	github.com/vardius/gorouter/v5 v5.0.0
	google.golang.org/grpc v1.85.0
)

require (
	golang.org/x/net v0.60.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/text v0.42.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260817212433-ac3dfec99bb1 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace github.com/vardius/gorouter/v5 => ../..
