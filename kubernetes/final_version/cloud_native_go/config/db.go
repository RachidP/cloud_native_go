package config

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
)

//mongodb_service: is the name of the service inside docker-compose
const (
	//mongodb_service = "mongodb-service"
	db_name         = "books"
	collection_name = "book"
)

var mongodb_service string
var mongodb_service_port string
var Collection *mongo.Collection

func init() {
	var err error
	var client *mongo.Client

	mongodb_service = "mongodb"
	mongodb_service_port := ":27017"
	//mongodb_service = "mongodb"
	//create a mongo.Client
	//http://mongodb:27017
	checkdbConnection()
	client, err = mongo.NewClient("mongodb://" + mongodb_service + mongodb_service_port)

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

func getMongoService() string {
	mongodb_service = os.Getenv("MONGODB_SERVICE_HOST")
	port := os.Getenv("MONGODB_SERVICE_PORT")
	if len(mongodb_service) == 0 {
		fmt.Println("NO services found")
		mongodb_service = "mongodb"

	}
	fmt.Printf("IP:%v   PORT:%v\n", mongodb_service, port)
	log.Printf("Rachid the address of the DB is: %v \n", mongodb_service)

	//return mongodb_service + ":" + port
	return mongodb_service
}

func checkdbConnection() {

	res, err := http.Get("http://mongodb:27017")
	if err != nil {
		fmt.Printf("Error during http.GET to mongo : %v", err)
		return
	}
	fmt.Println("You have a get from mongo db")
	fmt.Printf("The response from mongo db is :\n----> %v\n", res)

}
