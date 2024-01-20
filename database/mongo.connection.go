package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBInstance() *mongo.Client {
	url := os.Getenv("MONGO_URL")

	if url == "" {
		log.Fatal("fatal load mongo db url ")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(url))

	if err != nil {
		log.Fatal("could not connect to mongo db cuz of ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connect to mongo db successfully")

	return client
}

var Client *mongo.Client = DBInstance()

func Collection(client *mongo.Client, collection_name string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("cluster0").Collection(collection_name)

	return collection
}
