package main

import (
	"flag"
	"fmt"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/golang/glog"
	pb "github.com/kirk91/misc/go/grpc-gateway/echopb"
)

var (
	port = flag.Int("port", 9090, "listen port")
)

type echoService struct{}

func (s *echoService) Echo(ctx context.Context, in *pb.StringMessage) (*pb.StringMessage, error) {
	return in, nil
}

func run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	pb.RegisterEchoServiceServer(grpcServer, &echoService{})
	return grpcServer.Serve(lis)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
