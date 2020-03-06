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

	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "postgres"
		dbName   = "postgres"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)

	usersController := controllers.NewUser(psqlInfo)

	r.HandleFunc("/", static.HomeView.ServeHTTP).Methods("GET")
	r.HandleFunc("/contact", static.ContactView.ServeHTTP).Methods("GET")
	r.HandleFunc("/faq", static.FAQView.ServeHTTP).Methods("GET")
	r.HandleFunc("/sign-up", usersController.New).Methods("GET")
	r.HandleFunc("/sign-up", usersController.Create).Methods("POST")
	r.HandleFunc("/login", usersController.LoginView.ServeHTTP).Methods("GET")
	r.HandleFunc("/login", usersController.SignIn).Methods("POST")
	r.HandleFunc("/cookietest", usersController.CookieTest).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(static.NotFoundView.ServeHTTP)
	fmt.Println("Running on port 3000")
	http.ListenAndServe(":3000", r)
}
