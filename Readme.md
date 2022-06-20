# Crypto Currency - API
Project developed in Golang with gRPC

## Proto file
For to create new protos files, remember to delete files ``proto/service_grpc.pb.go`` ``proto/service.pb.go`` and run command below

> protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/service.proto``

## Access
Clone this repository and execute command ``go run main.go`` into project root

## Test
For a better view I suggest use of plugin [Mocha Test Explorer](https://marketplace.visualstudio.com/items?itemName=hbenl.vscode-mocha-test-adapter)

Or run the tests with the command ``go test ./...`` into project root

## Requests
For requests, I recommend use of client [BloomRPC](https://github.com/bloomrpc/bloomrpc)

_Remember to import the ``proto/service.proto`` file in your client_


___

###### Project developed by [Daniel Carvalho](https://www.linkedin.com/in/daniel-carvalho-7844b6107/)
