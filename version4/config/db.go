package config

import (
	"context"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
)

var Collection *mongo.Collection

func init() {
	var err error
	var client *mongo.Client

	//create a mongo.Client
	client, err = mongo.NewClient("mongodb://localhost:27017")
	if err != nil {
		log.Printf("could not create a mongo client:  %v", err)
		panic(err)
	}

	//connect it to your running MongoDB server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//Connect initializes the Client by starting background monitoring goroutines.
	err = client.Connect(ctx)
	if err != nil {
		log.Printf("could not connect with mongodb:  %v", err)
		panic(err)
	}

	//retreive the DB (named books) and the Collection (named book)
	Collection = client.Database("books").Collection("book")

}
