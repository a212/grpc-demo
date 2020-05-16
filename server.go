package main

import (
	"fmt"
	api "github.com/a212/grpc-demo/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"io/ioutil"
	"net"
	"net/http"
)

const GetBodyLen = 1024

func main() {
	listener, err := net.Listen("tcp", ":4242")

	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	server := Server{
		Results: make([]Result, 0),
	}

	api.RegisterDemoServer(grpcServer, &server)
	grpcServer.Serve(listener)
}

type Result struct {
	Headers string
	Body    []byte
	BodyPos int
}

type Server struct {
	Results []Result
}

func (s *Server) Do(c context.Context, request *api.Request) (*api.Handle, error) {
	resp, err := http.Get(request.Url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	handle := api.Handle{
		Id: uint64(len(s.Results)),
	}
	result := Result{
		Headers: fmt.Sprintf("%v", resp.Header),
		Body:    body,
		BodyPos: 0,
	}
	s.Results = append(s.Results, result)
	return &handle, nil
}

func (s *Server) GetHeaders(c context.Context, handle *api.Handle) (*api.Headers, error) {
	if handle.Id >= uint64(len(s.Results)) {
		return nil, fmt.Errorf("Invalid handle: %v", handle.Id)
	}
	headers := api.Headers{
		Content: s.Results[handle.Id].Headers,
	}
	return &headers, nil
}

func (s *Server) GetBody(c context.Context, handle *api.Handle) (*api.Body, error) {
	if handle.Id >= uint64(len(s.Results)) {
		return nil, fmt.Errorf("Invalid handle: %v", handle.Id)
	}
	r := &s.Results[handle.Id]
	sz := GetBodyLen
	if len(r.Body)-r.BodyPos-sz < 0 {
		sz = len(r.Body) - r.BodyPos
	}
	body := api.Body{
		Content: r.Body[r.BodyPos : r.BodyPos+sz],
	}
	r.BodyPos += sz
	return &body, nil
}
