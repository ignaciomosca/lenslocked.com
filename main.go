package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"lenslocked.com/controllers"
)

func main() {
	r := mux.NewRouter()
	static := controllers.NewStatic()
	usersController := controllers.NewUser()

	r.HandleFunc("/", static.HomeView.ServeHTTP).Methods("GET")
	r.HandleFunc("/contact", static.ContactView.ServeHTTP).Methods("GET")
	r.HandleFunc("/faq", static.FAQView.ServeHTTP).Methods("GET")
	r.HandleFunc("/sign-up", usersController.New).Methods("GET")
	r.HandleFunc("/sign-up", usersController.Create).Methods("POST")
	r.NotFoundHandler = http.HandlerFunc(static.NotFoundView.ServeHTTP)
	fmt.Println("Running on port 3000")
	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
