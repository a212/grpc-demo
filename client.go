package main

import (
	"context"
	"fmt"
	api "github.com/a212/grpc-demo/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"os"
)

func main() {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	args := os.Args
	if len(os.Args) != 2 {
		grpclog.Fatalf("Usage: %v url", args[0])
		os.Exit(1)
	}
	conn, err := grpc.Dial("127.0.0.1:4242", opts...)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}

	defer conn.Close()

	client := api.NewDemoClient(conn)
	request := api.Request{
		Url: args[1],
	}

	handle, err := client.Do(context.Background(), &request)

	if err != nil {
		grpclog.Fatalf("fail to do request: %v", err)
	}

	headers, err := client.GetHeaders(context.Background(), handle)

	if err != nil {
		grpclog.Fatalf("fail to get headers: %v", err)
	}

	fmt.Printf("Headers: %v\n", headers.Content)

	body, err := client.GetBody(context.Background(), handle)

	if err != nil {
		grpclog.Fatalf("fail to get body: %v", err)
	}
	fmt.Printf("Body: %v\n", string(body.Content[:100]))
}
