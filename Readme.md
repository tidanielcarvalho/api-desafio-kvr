# Crypto Currency - API
Project developed in Golang with gRPC

## Start Application
1. Clone this repository

2. Duplicate env.example file to .env

3. Start container, run ``docker-compose -f docker/docker-compose.yml up``

4. Run ``go run main.go`` into project root


> Log debug is enable, to debug disable change LOG_DEBUG variable in .env

## Database
The database is docker container with mongo image

To view the data access ``localhost:9000`` with user ``root`` and password ``root``

## Proto file
To create new protos files, remember to delete files ``proto/service_grpc.pb.go`` ``proto/service.pb.go`` and run command below

> protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/service.proto

## Test
For a better view I suggest use of plugin [Mocha Test Explorer](https://marketplace.visualstudio.com/items?itemName=hbenl.vscode-mocha-test-adapter)

Or run the tests with the command ``go test -v ./...`` into project root

## Requests
To requests, I recommend use of client [BloomRPC](https://github.com/bloomrpc/bloomrpc)

_Remember to import the ``proto/service.proto`` file in your client_

## Default Config
> - Ports:
>
>   - Application: 55555
>
>   - Container
>
>       - Mongo: 27017->27017
>
>       - Mongo-Express: 8081->9000

___
## Problems?
 Send mail to ti.danielcarvalho@gmail.com or [LinkedIn](https://www.linkedin.com/in/daniel-carvalho-7844b6107/)


###### Project developed by [Daniel Carvalho](https://www.linkedin.com/in/daniel-carvalho-7844b6107/)
