package books

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Book type with Name, Author and Isbn
type Book struct {
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Isbn   string  `json:"Isbn"`
	Price  float32 `json:"price,omitempty"`
}

var books = map[string]Book{
	"0345391802":  Book{Title: "The Hitchhiker's Guide to the Galaxy", Author: "Douglas Adams", Isbn: "0345391802"},
	"0000000000":  Book{Title: "Cloud Native Go", Author: "M.-Leander Reimer", Isbn: "0000000000"},
	"00000000045": Book{Title: "Mastering Kubernetes", Author: "Gigi Sayfan", Isbn: "00000000045", Price: 20.0},
}

// AllBooks returns a slice of all books
func AllBooks() ([]Book, error) {
	values := make([]Book, len(books))
	idx := 0
	for _, book := range books {
		values[idx] = book
		idx++
	}
	return values, nil
}

// OneBook search a book by isbn
func OneBook(r *http.Request) (Book, error) {
	bk := Book{}
	isbn := r.FormValue("isbn")
	if isbn == "" {
		return bk, errors.New("400. Bad Request.")
	}

	bk, found := books[isbn]
	if !found {
		return bk, fmt.Errorf("could not find the book (%v)in the DB", isbn)
	}

	return bk, nil
}

func PutBook(r *http.Request) (Book, error) {
	// get form values
	bk := Book{}
	bk.Isbn = r.FormValue("isbn")
	bk.Title = r.FormValue("title")
	bk.Author = r.FormValue("author")
	p := r.FormValue("price")

	// validate form values
	if bk.Isbn == "" || bk.Title == "" || bk.Author == "" || p == "" {
		return bk, errors.New("400. Bad request. All fields must be complete.")
	}

	// convert form values
	f64, err := strconv.ParseFloat(p, 32)
	if err != nil {
		return bk, errors.New("406. Not Acceptable. Price must be a number.")
	}
	bk.Price = float32(f64)

	err = addBook(bk)
	if err != nil {
		return bk, errors.New("500. Internal Server Error." + err.Error())
	}
	return bk, nil
}

// CreateBook creates a new Book if it does not exist
func addBook(book Book) error {
	//check if exist
	_, exists := books[book.Isbn]
	if exists {
		fmt.Errorf("book already exist")
	}
	books[book.Isbn] = book
	return nil
}

func UpdateBook(r *http.Request) (Book, error) {
	// get form values
	bk := Book{}
	bk.Isbn = r.FormValue("isbn")
	bk.Title = r.FormValue("title")
	bk.Author = r.FormValue("author")
	p := r.FormValue("price")

	if bk.Isbn == "" || bk.Title == "" || bk.Author == "" || p == "" {
		return bk, errors.New("400. Bad Request. Fields can't be empty.")
	}

	// convert form values
	f64, err := strconv.ParseFloat(p, 32)
	if err != nil {
		return bk, errors.New("406. Not Acceptable. Enter number for price.")
	}
	bk.Price = float32(f64)

	//check if exit in the DB
	_, exists := books[bk.Isbn]
	if !exists {
		return bk, fmt.Errorf("book does not exist!")

	}
	// insert values
	books[bk.Isbn] = bk

	return bk, nil
}

func DeleteBook(r *http.Request) error {
	isbn := r.FormValue("isbn")
	_, exist := books[isbn]
	if isbn == "" || !exist {
		log.Printf("isbn=%v not valid\n", isbn)
		return errors.New("400. Bad Request.")
	}

	// _, err := config.DB.Exec("DELETE FROM books WHERE isbn=$1;", isbn)
	// if err != nil {
	// 	return errors.New("500. Internal Server Error")
	// }

	delete(books, isbn)
	return nil
}
