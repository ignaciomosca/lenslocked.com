package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

type User struct {
	Name string
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, err := template.ParseFiles("templates/home.gohtml")
	if err != nil {
		panic(err)
	}
	data := User{Name: "Juan"}
	t.Execute(w, data)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, err := template.ParseFiles("templates/contact.gohtml")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, err := template.ParseFiles("templates/faq.gohtml")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, err := template.ParseFiles("templates/notfound.gohtml")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

func main() {
	r := mux.NewRouter()
	template.New("blah")
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/faq", faq)
	r.NotFoundHandler = http.HandlerFunc(notFound)
	fmt.Println("Running on port 3000")
	http.ListenAndServe(":3000", r)
}
