package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func ConnectMongoDB() *mongo.Database {
	clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB...")
	database := client.Database("shopping")
	return database
}
