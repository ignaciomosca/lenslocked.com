package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"lenslocked.com/views"
	"net/http"
)

var (
	homeTemplate     *views.View
	contactTemplate  *views.View
	faqTemplate      *views.View
	notFoundTemplate *views.View
)

type User struct {
	Name string
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	data := User{Name: "Juan"}
	if err := homeTemplate.Template.Execute(w, data); err != nil {
		panic(err)
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := contactTemplate.Template.Execute(w, nil); err != nil {
		panic(err)
	}
}

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := faqTemplate.Template.Execute(w, nil); err != nil {
		panic(err)
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	notFoundTemplate.Template.Execute(w, nil)
}

func main() {
	r := mux.NewRouter()
	homeTemplate = views.NewFiles("views/home.gohtml")
	contactTemplate = views.NewFiles("views/contact.gohtml")
	faqTemplate = views.NewFiles("views/faq.gohtml")
	contactTemplate = views.NewFiles("views/contact.gohtml")
	notFoundTemplate = views.NewFiles("views/notfound.gohtml")

	template.New("blah")
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/faq", faq)
	r.NotFoundHandler = http.HandlerFunc(notFound)
	fmt.Println("Running on port 3000")
	http.ListenAndServe(":3000", r)
}
