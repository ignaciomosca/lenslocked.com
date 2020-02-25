package controllers

import "lenslocked.com/views"

import "net/http"

import "fmt"

import "github.com/gorilla/schema"

func NewUser() *Users {
	return &Users{
		NewView: views.NewFiles("bootstrap", "views/users/new.gohtml"),
	}
}

type Users struct {
	NewView *views.View
}

type SignupForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	user := new(SignupForm)
	decoder := schema.NewDecoder()
	if err := decoder.Decode(user, r.PostForm); err != nil {
		panic(err)
	}

	fmt.Println("email and password", user.Email, user.Password)
	u.NewView.Render(w, nil)
}
