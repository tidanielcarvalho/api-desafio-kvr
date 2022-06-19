package main

import (
	"api-desafio-kvr/controllers"
	"api-desafio-kvr/helpers"
	"api-desafio-kvr/proto"
	"api-desafio-kvr/repositories/migration"
	"api-desafio-kvr/repositories/mongodb"
	"net"

	"google.golang.org/grpc"
)

var logger = &helpers.Log{}

func main() {
	client, ctx, cancel, _ := mongodb.Connect()

	collection := mongodb.GetDataBase(client)
	app := &controllers.AppServer{Database: collection}

	migration.CreateInitialCryptosBulk(app.Database)
	controllers.InitializeChanToStream()
	InitializeGRPC(app)

	mongodb.Disconnect(client, ctx, cancel)
}

func InitializeGRPC(app *controllers.AppServer) {
	logger.Info("", "Initializing gRPC service")

	grpc := grpc.NewServer()
	proto.RegisterEndPointCryptosServer(grpc, app)

	port := "5000"
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err.Error)
	}
	logger.Info("", "gRPC service running on the port "+port)
	err = grpc.Serve(listener)
	if err != nil {
		panic(err.Error)
	}
}
