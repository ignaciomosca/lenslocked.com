package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello World Chumbawamba</h1>")
	} else if r.URL.Path == "/contact" {
		fmt.Fprint(w, "<h1>Mandale un mail a Carlitos</h1>")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>404: Not Found</h1>")
	}
}

func main() {
	http.HandleFunc("/", handlerFunc)
	fmt.Println("Running on port 3000")
	http.ListenAndServe(":3000", nil)
}
