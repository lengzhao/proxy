# proxy

## build

```bash
go build -o client.bin ./client/main.go
go build -o server.bin ./server/main.go
go build -o proxy.bin main.go
```

## run

1. client -> server
   1. run server(terminal1): ./server.bin
   2. run client(terminal2): ./client.bin

```bash
% ./server.bin 
server listening at [::]:50051
Received: world

% ./client.bin
Greeting: Hello world
```

1. client -> proxy -> server
   1. run server(terminal1): ./server.bin
   2. run proxy(terminal2): ./proxy.bin
   3. run client(terminal3): ./client.bin -addr 127.0.0.1:50052

```bash
% ./server.bin 
server listening at [::]:50051
Received: world

% ./proxy.bin 
2022/09/04 11:51:44 server listening at [::]:50052
proxy: /helloworld.Greeter/SayHello

% ./client.bin -addr 127.0.0.1:50052
Greeting: Hello world
```
