package mongodb

import (
	"api-desafio-kvr/helpers"
	"context"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var logger = &helpers.Log{}

func getURI() string {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("MONGODB", "", err)
	}
	// Cluster
	// URI := "mongodb+srv://" + os.Getenv("MONGODB_USER") + ":" + os.Getenv("MONGODB_PASS") + "@clusterapi.t6vp0.mongodb.net/?retryWrites=true&w=majority"

	// Container
	URI := "mongodb://root:root@127.0.0.1:27017/?authSource=admin"

	return URI
}

func Connect() (*mongo.Client, context.Context, context.CancelFunc, error) {
	logger.Info("", "Starting database connection")
	URI := getURI()

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(URI).
		SetServerAPIOptions(serverAPIOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		logger.Fatal("", err.Error(), err)
	}

	return client, ctx, cancel, nil
}

func Disconnect(client *mongo.Client, context context.Context, cancel context.CancelFunc) {
	defer cancel()

	err := client.Disconnect(context)
	if err != nil {
		logger.Fatal("", err.Error(), err)
	}

	logger.Info("", "Closing database connection")
}
