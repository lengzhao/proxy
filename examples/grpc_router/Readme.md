# router

## build

```bash
go build -o client.bin ../grpc_hello/client/main.go
go build -o server.bin ../grpc_hello/server/main.go
go build -o router.bin main.go
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

1. config:
   1. vim conf.json
2. client -> router -> server
   1. run server(terminal1): ./server.bin
   2. run router(terminal2): ./router.bin
   3. run client(terminal3): ./client.bin -addr 127.0.0.1:50052

```bash
% ./server.bin 
server listening at [::]:50051
Received: world

% ./router.bin 
2022/09/04 11:51:44 server listening at [::]:50052
route: /helloworld.Greeter/SayHello 127.0.0.1:50051

% ./client.bin -addr 127.0.0.1:50052
Greeting: Hello world
```
