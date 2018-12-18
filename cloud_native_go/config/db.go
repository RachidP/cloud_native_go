package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
)

//mongodb_service: is the name of the service inside docker-compose
const (
	mongodb_service = "mongodb_service"
	db_name         = "books"
	collection_name = "book"
)

var Collection *mongo.Collection

func init() {
	var err error
	var client *mongo.Client

	//create a mongo.Client
	client, err = mongo.NewClient("mongodb://" + mongodb_service + ":27017")

	if err != nil {
		log.Printf("could not create a mongo client:  %v", err)
		panic(err)
	}
	//
	//connect it to your running MongoDB server
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	//Connect initializes the Client by starting background monitoring goroutines.
	err = client.Connect(ctx1)
	if err != nil {
		log.Printf("could not connect with mongodb:  %v", err)
		panic(err)
	}

	// verifies if the client can connect to the monogoDB
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()
	err = client.Ping(ctx2, nil)
	if err != nil {
		log.Printf("No ping to mongo db! could not connect to mongodb:  %v", err)
		panic(err)
	}

	//retreive the DB (named books) and the Collection (named book)
	Collection = client.Database(db_name).Collection(collection_name)
	fmt.Println("Mongo DB is connected and is ready for query!")

}
