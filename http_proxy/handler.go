package http_proxy

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
)

// get address from request
func AutoHandleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}
	var method, host, address string
	fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &host)
	hostPortURL, err := url.Parse(host)
	if err != nil {
		log.Println(err)
		return
	}

	if hostPortURL.Opaque == "443" { //https
		address = hostPortURL.Scheme + ":443"
	} else { //http
		if !strings.Contains(hostPortURL.Host, ":") {
			address = hostPortURL.Host + ":80"
		} else {
			address = hostPortURL.Host
		}
	}

	server, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(address, err)
		return
	}
	log.Println("proxy:", address)
	if method == "CONNECT" {
		fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		server.Write(b[:n])
	}

	c1 := make(chan int)
	c2 := make(chan int)

	go func() {
		io.Copy(server, client)
		c1 <- 1
	}()

	go func() {
		io.Copy(client, server)
		c2 <- 1
	}()

	for i := 0; i < 2; i++ {
		select {
		case <-c1:
			server.Close()
			client.Close()
		case <-c2:
			server.Close()
			client.Close()
		}
	}
}

// force proxy to next hop
func ProxyToNextHop(next string, client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	server, err := net.Dial("tcp", next)
	if err != nil {
		log.Println("ProxyToNextHop", next, err)
		return
	}
	c1 := make(chan int)
	c2 := make(chan int)

	go func() {
		io.Copy(server, client)
		c1 <- 1
	}()

	go func() {
		io.Copy(client, server)
		c2 <- 1
	}()

	for i := 0; i < 2; i++ {
		select {
		case <-c1:
			server.Close()
			client.Close()
		case <-c2:
			server.Close()
			client.Close()
		}
	}
}
