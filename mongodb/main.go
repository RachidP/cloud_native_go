package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

var collection *mongo.Collection

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
	//retreive the db (named books) and the collection (named book)
	collection = client.Database("books").Collection("book")

}

func main() {
	insertOneDocument()
	getNumDocument()
	getAllDocument()
	dropDocument()
	getNumDocument()
	updateDocument()
	getAllDocument()
}

type Book struct {
	Title  string  `json:"title" bson:"title"`
	Author string  `json:"author" bson:"author"`
	Isbn   int     `json:"Isbn" bson:"isbn"`
	Price  float64 `json:"price,omitempty"  bson:"price"`
}

//insertOneDocument add a document in a collection
//in this version I use the bson.Marshal() to tranform
// golang data struct to bson document
func insertOneDocument() {

	b := Book{
		Title:  "lalaladfdfdfla",
		Author: "monononon",
		Isbn:   75565843,
		Price:  84,
	}

	bsonDoc, err := bson.Marshal(b)
	if err != nil {
		fmt.Printf("could not marshal the data %v", err)
		panic(err)
	}

	res, err := collection.InsertOne(context.Background(), bsonDoc)

	//use with bson.D{}
	//res, err := collection.InsertOne(context.Background(), bson.D{{"title", b.Title}, {"author", b.Author}, {"isbn", b.Isbn}, {"price", b.Price}})
	if err != nil {
		log.Printf("could not insert data into collection:  %v\n", err)
		panic(err)
	}
	fmt.Printf("res:  %v\n", res)

}

//get the number of documents in a collection
func getNumDocument() {
	num, err := collection.CountDocuments(context.Background(), nil)
	if err != nil {
		log.Printf("could not count the number of documents:  %v\n", err)
		panic(err)
	}
	fmt.Printf("The are %v documents in your collection\n", num)

}

// find all documents
//in this solution I use Decode to decode from BSON to golang data struct
func getAllDocument() {

	var b Book

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(context.Background()) {

		cur.Decode(&b)
		if err != nil {
			fmt.Printf("could not decode the data: %v ", err)
			log.Fatal(err)
		}
		fmt.Printf("book: %#v\n", b)
		fmt.Printf("isbn: %v\n", b.Isbn)

	}
	if err := cur.Err(); err != nil {
		log.Printf("could not find data:  %v", err)
		panic(err)
	}
}

//drop a document
func dropDocument() {
	filter := bson.M{"isbn": 4543534}
	res, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Printf("could not delete the isbn:Internal Error %v\n", err)
		panic(err)
	}
	//no query match
	if res.DeletedCount == 0 {

		fmt.Printf("Isbn does not exist in db. number of item deleted = %v \n", res.DeletedCount)
	}

	fmt.Printf("document deleted\n")
}

//update a document
func updateDocument() {
	//filter := bson.M{"isbn": 88889}
	filter := bson.M{"isbn": 4543534}

	newData := bson.M{"$set": bson.M{"isbn": 1, "author": "updated author", "title": "updated title", "price": 20}}
	res, err := collection.UpdateOne(context.Background(), filter, newData)

	if err != nil {
		log.Printf(" could not update data (Internal Error): %v\n", err)
		panic(err)
	}
	if res.MatchedCount == 0 {
		fmt.Printf("Isbn does not exist in db. number of item deleted = %v \n", res.ModifiedCount)
	}

	log.Printf("Data updated successfully\n")
}

/*******************************************************************************/
// find all documents
//in this solution I use Decode to decode from BSON to bson.M data type
//after that it you can get the value by key. examlpe res["author"].
func getAllDocument2() {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(context.Background()) {
		var res bson.M
		err = cur.Decode(&res)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(res)
		fmt.Println(res["author"])

	}
	if err := cur.Err(); err != nil {
		log.Printf("could not find data:  %v", err)
		panic(err)
	}
}

// find all documents
//in this solution I use Decode using DecodeBytes to decode from BSON to bson.Raw
//after that use lookup method to access to data
func getAllDocument3() {

	cur, err := collection.Find(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		raw, err := cur.DecodeBytes()
		if err != nil {
			log.Fatal(err)
		}
		//with lookup
		//	fmt.Println("raw", raw)

		fmt.Printf("title= %v\n", raw.Lookup("title"))
		fmt.Printf("author=%v\n", raw.Lookup("author"))
		fmt.Printf("isbn=%v\n", raw.Lookup("isbn").Int64())
		fmt.Printf("price=%v\n", raw.Lookup("price").Double())

	}
	if err := cur.Err(); err != nil {
		log.Printf("could not find data:  %v", err)
		panic(err)
	}
}

//insertOneDocument add a document in a collection
//in this version 2 I use the bson.M{} or bson.D{} for trasforming
//golang data into bson data
func insertOneDocument2() {

	b := Book{
		Title:  "tu si",
		Author: "Dante Alleghieri",
		Isbn:   757843,
		Price:  14,
	}

	//use with bson.M{}
	//insert data to DB
	res, err := collection.InsertOne(context.Background(), bson.M{"title": b.Title, "author": b.Author, "isbn": b.Isbn, "price": b.Price})

	//use with bson.D{}
	//res, err := collection.InsertOne(context.Background(), bson.D{{"title", b.Title}, {"author", b.Author}, {"isbn", b.Isbn}, {"price", b.Price}})

	if err != nil {
		log.Printf("could not insert data into collection:  %v\n", err)
		panic(err)
	}
	fmt.Printf("res:  %v\n", res)

}
