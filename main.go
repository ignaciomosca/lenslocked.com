package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

var homeTemplate *template.Template

type User struct {
	Name string
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	data := User{Name: "Juan"}
	if err := homeTemplate.Execute(w, data); err != nil {
		panic(err)
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, err := template.ParseFiles("views/contact.gohtml")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, err := template.ParseFiles("views/faq.gohtml")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, err := template.ParseFiles("views/notfound.gohtml")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

func main() {
	r := mux.NewRouter()
	var err error
	homeTemplate, err = template.ParseFiles("views/home.gohtml")
	if err != nil {
		panic(err)
	}
	template.New("blah")
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/faq", faq)
	r.NotFoundHandler = http.HandlerFunc(notFound)
	fmt.Println("Running on port 3000")
	http.ListenAndServe(":3000", r)
}
