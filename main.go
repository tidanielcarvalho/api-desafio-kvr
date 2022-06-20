package main

import (
	"api-desafio-kvr/controllers"
	"api-desafio-kvr/helpers"
	"api-desafio-kvr/proto"
	"api-desafio-kvr/repositories/migration"
	"api-desafio-kvr/repositories/mongodb"
	"net"
	"os"

	"google.golang.org/grpc"
)

var logger = &helpers.Log{}

func main() {
	logger.Info("", "Starting services to application")
	client, ctx, cancel, _ := mongodb.Connect()

	collection := mongodb.GetDataBase(client)
	app := &controllers.AppServer{Database: collection}

	migration.CreateInitialCryptosBulk(app.Database)

	controllers.StartChanToStream()
	StartGRPC(app)

	mongodb.Disconnect(client, ctx, cancel)
}

func StartGRPC(app *controllers.AppServer) {
	logger.Info("", "Starting gRPC service")

	grpc := grpc.NewServer()
	proto.RegisterEndPointCryptosServer(grpc, app)

	port := os.Getenv("PORT")
	if port == "" {
		logger.Warn("", "Env PORT is empty or not found")
		port = "55555"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Fatal("", err.Error(), err)
	}

	logger.Info("", "gRPC service running on the port "+port)
	err = grpc.Serve(listener)
	if err != nil {
		logger.Fatal("", err.Error(), err)
	}
}
