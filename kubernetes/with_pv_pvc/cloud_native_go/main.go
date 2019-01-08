package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RachidP/cloud_native_go/books"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/books", books.Index)
	http.HandleFunc("/books/show", books.Show)
	http.HandleFunc("/books/create", books.Create)
	http.HandleFunc("/books/create/process", books.CreateProcess)
	http.HandleFunc("/books/update", books.Update)
	http.HandleFunc("/books/update/process", books.UpdateProcess)
	http.HandleFunc("/books/delete/process", books.DeleteProcess)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.ListenAndServe(port(), nil)
	log.Println("The site is running!")
}

func port() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "9090"
	}
	log.Println("Connected to port: ", port)

	return ":" + port
}

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/books", http.StatusSeeOther)

}
