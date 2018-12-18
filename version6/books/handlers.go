package books

import (
	"net/http"

	"github.com/RachidP/exercises/cloud_native_go/version6/config"
)

//Index  handle the index page.
func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	bks, err := AllBooks()
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	config.TPL.ExecuteTemplate(w, "books.html", bks)
}

//Show show the page of a single book
func Show(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	bk, err := serveOneBook(r)

	switch {
	case err.Error() == "book not found":
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	config.TPL.ExecuteTemplate(w, "show.html", bk)
}

//Create handle request for showing the page for creating a book.
func Create(w http.ResponseWriter, r *http.Request) {
	config.TPL.ExecuteTemplate(w, "create.html", nil)
}

//CreateProcess Handle the creating book process.
func CreateProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	bk, err := AddBook(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		//		http.Error(w, http.StatusText(406), http.StatusNotAcceptable)
		return
	}

	config.TPL.ExecuteTemplate(w, "created.html", bk)
}

//Update handle the request for showing a update page.
func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	bk, err := serveOneBook(r)
	// switch {
	// case err == sql.ErrNoRows:
	// 	http.NotFound(w, r)
	// 	return
	// case err != nil:
	// 	http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	// 	return
	// }
	if err != nil {

		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	config.TPL.ExecuteTemplate(w, "update.html", bk)
}

//UpdateProcess handle the process for updating a book.
func UpdateProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	bk, err := UpdateBook(r)
	if err != nil {
		//http.Error(w, http.StatusText(406), http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	config.TPL.ExecuteTemplate(w, "updated.html", bk)
}

//DeleteProcess handle a delete book process
func DeleteProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	_, err := DeleteBook(r)
	if err != nil {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/books", http.StatusSeeOther)
}
