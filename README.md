# proxy

## grpc proxy

1. examples/grpc_hello
2. examples/grpc_router
   1. go install github.com/lengzhao/proxy/examples/grpc_router@latest
   1. add conf.json in dir

    ```json
    {
        "route":{
            "helloworld.Greeter":"127.0.0.1:50051"
        }
    }
    ```

    1. run router: grpc_router
