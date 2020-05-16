# rpc-demo

## Building

		protoc api.proto --go_out=plugins=grpc:.
		go build client.go
		go build server.go

## Running

		server &
		client http://example.com