package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"

	"github.com/lengzhao/proxy/grpc_proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Config struct {
	Route map[string]string `json:"route,omitempty"`
}

func main() {
	port := flag.Int("port", 50052, "The server port")
	confFile := flag.String("c", "conf.json", "config file")
	flag.Parse()
	var conf Config
	data, err := ioutil.ReadFile(*confFile)
	if err != nil {
		log.Fatalln("fail to open config file:", err)
	}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Fatalln("error json file:", err)
	}
	if len(conf.Route) == 0 {
		log.Fatalln("not found any route item")
	}

	director := func(ctx context.Context, fullMethodName string) (*grpc.ClientConn, error) {
		arrs := strings.Split(fullMethodName, "/")
		if len(arrs) < 2 {
			return nil, status.Error(codes.InvalidArgument, fullMethodName)
		}
		svcName := arrs[1]
		endpoint := conf.Route[svcName]
		if len(endpoint) == 0 {
			log.Println("not found endpoint:", fullMethodName)
			return nil, status.Error(codes.Unauthenticated, "")
		}
		log.Println("route:", fullMethodName, endpoint)
		return grpc.DialContext(ctx, endpoint,
			grpc.WithCodec(grpc_proxy.ProxyCodec{}),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	ops := grpc_proxy.GetServerOptions(director)
	server := grpc.NewServer(ops...)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
