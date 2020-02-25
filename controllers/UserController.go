package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/views"
)

func NewUser() *Users {
	return &Users{
		NewView: views.NewFiles("bootstrap", "users/new"),
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
	if err := parse(r, user); err != nil {
		panic(err)
	}

	fmt.Println("email and password", user.Email, user.Password)
	u.NewView.Render(w, nil)
}
