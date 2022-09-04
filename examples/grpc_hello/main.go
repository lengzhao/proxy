package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/lengzhao/proxy/grpc_proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	port := flag.Int("port", 50052, "The server port")
	endpoint := flag.String("endpoint", "127.0.0.1:50051", "The grpc server address")
	flag.Parse()
	director := func(ctx context.Context, fullMethodName string) (*grpc.ClientConn, error) {
		fmt.Println("proxy:", fullMethodName)
		return grpc.DialContext(ctx, *endpoint,
			grpc.WithCodec(grpc_proxy.ProxyCodec{}),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	server := grpc.NewServer(
		grpc.CustomCodec(grpc_proxy.ProxyCodec{}),
		grpc_proxy.GetServerOption(director))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
