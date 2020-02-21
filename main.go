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

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := homeTemplate.Template.ExecuteTemplate(w, homeTemplate.Layout, nil)
	if err != nil {
		panic(err)
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := contactTemplate.Template.ExecuteTemplate(w, contactTemplate.Layout, nil); err != nil {
		panic(err)
	}
}

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := faqTemplate.Template.ExecuteTemplate(w, faqTemplate.Layout, nil); err != nil {
		panic(err)
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	notFoundTemplate.Template.ExecuteTemplate(w, notFoundTemplate.Layout, nil)
}

func main() {
	r := mux.NewRouter()
	homeTemplate = views.NewFiles("bootstrap", "views/home.gohtml")
	contactTemplate = views.NewFiles("bootstrap", "views/contact.gohtml")
	faqTemplate = views.NewFiles("bootstrap", "views/faq.gohtml")
	contactTemplate = views.NewFiles("bootstrap", "views/contact.gohtml")
	notFoundTemplate = views.NewFiles("bootstrap", "views/notfound.gohtml")

	template.New("blah")
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/faq", faq)
	r.NotFoundHandler = http.HandlerFunc(notFound)
	fmt.Println("Running on port 3000")
	http.ListenAndServe(":3000", r)
}
