package books

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/RachidP/exercises/cloud_native_go/version4/config"
	"github.com/mongodb/mongo-go-driver/bson"
)

// Book type with Name, Author and Isbn
type Book struct {
	Title  string  `json:"title" bson:"title"`
	Author string  `json:"author" bson:"author"`
	Isbn   int     `json:"Isbn" bson:"isbn"`
	Price  float64 `json:"price,omitempty" bson:"price"`
}

// AllBooks get the list of the books from DB
func AllBooks() ([]Book, error) {

	books := make([]Book, 0)

	var b Book

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//get result of the query as a cursor
	cur, err := config.Collection.Find(ctx, nil)
	if err != nil {
		log.Printf("could not query the collection: %v \n", err)
		return nil, errors.New("500. Internal Server Error." + err.Error())

	}
	defer cur.Close(ctx)

	//iterate over the cursor
	for cur.Next(context.Background()) {

		//decode the data
		cur.Decode(&b)
		if err != nil {
			log.Printf("could not decode the data: %v \n", err)
			return nil, errors.New("500. Internal Server Error." + err.Error())

		}

		books = append(books, b)

	}
	if err := cur.Err(); err != nil {
		log.Printf("could not iterate with the cursor: %v \n", err)
		return nil, errors.New("500. Internal Server Error." + err.Error())
	}

	//	fmt.Printf("list of all book: %#v\n", books)

	return books, nil
}

//AddBook Make the operation for adding a book in the db
func AddBook(r *http.Request) (Book, error) {
	//parse request
	var noBook Book
	bk, err := parseRequest(r)
	if err != nil {
		return bk, err
	}

	//check if the book already exist in the DB
	tmpBook, err := getBookFromDB(bk.Isbn)
	if err != nil && err.Error() != "Book not found" {

		return bk, err
	}

	//if the book already exist
	if tmpBook != noBook {
		log.Printf("book already exist!\n ")
		return bk, fmt.Errorf("406. Not Acceptable. book with Isbn (%v) already exist.", bk.Isbn)

	}

	//Add the book in the DB
	bsonDoc, err := bson.Marshal(bk)
	if err != nil {
		log.Printf("could not marshal the data %v \n", err)
		return bk, errors.New("500. Internal Server Error." + err.Error())

	}
	_, err = config.Collection.InsertOne(context.Background(), bsonDoc)
	if err != nil {
		log.Printf("could not insert data into collection:  %v\n", err)
		return bk, errors.New("500. Internal Server Error." + err.Error())

	}

	fmt.Printf("Book Added to DB! \n")
	return bk, nil
}

//parseRequest check if the data received are correct and return the book.
func parseRequest(r *http.Request) (Book, error) {

	var bk Book
	var err error

	price := r.FormValue("price")

	// validate form values
	if r.FormValue("isbn") == "" || r.FormValue("title") == "" || r.FormValue("author") == "" || price == "" {
		return bk, errors.New("400. Bad request. All fields must be complete.")
	}

	// convert form values
	bk.Price, err = strconv.ParseFloat(price, 64)
	if err != nil {
		log.Printf("could not convert string to float: %v\n ", err)
		return bk, errors.New("406. Not Acceptable. Price must be a number.")
	}

	bk.Isbn, err = strconv.Atoi(r.FormValue("isbn"))
	if err != nil {
		log.Printf("could not convert string to int: %v\n ", err)
		return bk, errors.New("406. Not Acceptable. Isbn must be a number.")
	}

	bk.Title = r.FormValue("title")
	bk.Author = r.FormValue("author")

	return bk, nil

}

//UpdateBook update the Book on the DB
func UpdateBook(r *http.Request) (Book, error) {
	var noBook Book
	bk, err := parseRequest(r)
	if err != nil {
		return bk, err
	}

	//oldisbn if the isbn has been update too
	old_isbn, err := strconv.Atoi(r.FormValue("old_isbn"))
	if err != nil {
		log.Printf("could not convert string to int: %v\n ", err)
		return bk, errors.New("406. Not Acceptable. Isbn must be a number.")
	}

	//if the isbn has been changed
	if old_isbn != bk.Isbn {
		//check if the isbn already exist
		tmpBook, err := getBookFromDB(bk.Isbn)
		if err != nil {

			return bk, err
		}

		if tmpBook != noBook {
			log.Printf("book already exist!\n ")
			return bk, fmt.Errorf("406. Not Acceptable. book with Isbn (%v) already exist.", bk.Isbn)
		}

	}

	bsonDoc := bson.M{"$set": bk}

	filter := bson.M{"isbn": old_isbn}

	res, err := config.Collection.UpdateOne(context.Background(), filter, bsonDoc)

	if err != nil {
		log.Printf(" could not update data (Internal Error): %v\n", err)
		return bk, errors.New("500. Internal Server Error." + err.Error())
	}

	if res.MatchedCount == 0 {
		fmt.Printf("Isbn does not exist in db. number of item deleted = %v \n", res.ModifiedCount)
		return bk, errors.New("400. Bad Request.")

	}

	log.Printf("Book updated!\n")

	return bk, nil
}

//serveOneBook check if the isbn is correct and return one book from DB
func serveOneBook(r *http.Request) (Book, error) {

	var bk Book
	if r.FormValue("isbn") == "" {
		return bk, errors.New("400. Bad Request.")
	}

	isbn, err := strconv.Atoi(r.FormValue("isbn"))
	if err != nil {
		log.Printf("could not convert string to int: %v\n ", err)
		return bk, errors.New("406. Not Acceptable. Isbn must be a number.")
	}

	return getBookFromDB(isbn)

}

//getBook return a book if exist!
func getBookFromDB(isbn int) (Book, error) {

	var bk Book
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"isbn": isbn}
	err := config.Collection.FindOne(ctx, filter).Decode(&bk)

	if err != nil {

		log.Printf("no book founded: %v \n", err)
		return bk, fmt.Errorf("Book not found")

	}

	fmt.Println(bk)
	return bk, nil
}

func DeleteBook(r *http.Request) (Book, error) {

	var bk Book
	if r.FormValue("isbn") == "" {
		return bk, errors.New("400. Bad Request.")
	}

	isbn, err := strconv.Atoi(r.FormValue("isbn"))
	if err != nil {
		log.Printf("could not convert string to int: %v\n ", err)
		return bk, errors.New("406. Not Acceptable. Isbn must be a number.")
	}

	filter := bson.M{"isbn": isbn}
	res, err := config.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Printf("could not delete the isbn:Internal Error %v\n", err)
		return bk, errors.New("500. Internal Server Error." + err.Error())

	}
	//no query match
	if res.DeletedCount == 0 {
		fmt.Printf("Isbn does not exist in db. number of item deleted = %v \n", res.DeletedCount)
		return bk, errors.New("400. Bad Request.")
	}

	fmt.Printf("book deleted\n")

	return bk, nil
}
