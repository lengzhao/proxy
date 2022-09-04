package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/lengzhao/proxy/http_proxy"
)

func main() {
	port := flag.Int("port", 8081, "The server port")
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	addr := fmt.Sprintf(":%d", *port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(err)
	}

	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}

		go http_proxy.AutoHandleClientRequest(client)
	}
}
