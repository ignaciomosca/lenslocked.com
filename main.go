package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"lenslocked.com/controllers"
	"lenslocked.com/views"
)

var (
	homeTemplate     *views.View
	contactTemplate  *views.View
	faqTemplate      *views.View
	notFoundTemplate *views.View
	signInTemplate   *views.View
	signUpTemplate   *views.View
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(homeTemplate.Render(w, nil))
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(contactTemplate.Render(w, nil))
}

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(faqTemplate.Render(w, nil))
}

func signIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(signInTemplate.Render(w, nil))
}

func signUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(signUpTemplate.Render(w, nil))
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
	signInTemplate = views.NewFiles("bootstrap", "views/sign-in.gohtml")
	signUpTemplate = views.NewFiles("bootstrap", "views/sign-up.gohtml")

	usersController := controllers.NewUser()

	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/contact", contact).Methods("GET")
	r.HandleFunc("/faq", faq).Methods("GET")
	r.HandleFunc("/sign-in", signIn).Methods("GET")
	r.HandleFunc("/sign-up", usersController.New).Methods("GET")
	r.HandleFunc("/sign-up", usersController.Create).Methods("POST")
	r.NotFoundHandler = http.HandlerFunc(notFound)
	fmt.Println("Running on port 3000")
	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
