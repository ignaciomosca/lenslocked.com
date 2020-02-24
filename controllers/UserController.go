package controllers

import "lenslocked.com/views"

import "net/http"

import "fmt"

func NewUser() *User {
	return &User{
		NewView: views.NewFiles("bootstrap", "views/users/new.gohtml"),
	}
}

type User struct {
	NewView *views.View
}

func (u *User) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}

func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	fmt.Println("email and password", email, password)
	u.NewView.Render(w, nil)
}
