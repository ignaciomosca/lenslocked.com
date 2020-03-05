package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/models"
	"lenslocked.com/views"
)

func NewUser(connectionInfo string) *Users {
	us, err := models.NewUserService(connectionInfo)
	if err != nil {
		panic(err)
	}

	return &Users{
		NewView:     views.NewFiles("bootstrap", "users/new"),
		LoginView:   views.NewFiles("bootstrap", "users/login"),
		UserService: us,
	}
}

type Users struct {
	NewView     *views.View
	LoginView   *views.View
	UserService *models.UserService
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}

func (u *Users) SignIn(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	login := new(LoginForm)
	if err := parse(r, login); err != nil {
		panic(err)
	}
	loggedUser, loginErr := u.UserService.Login(login.Email, login.Password)
	if loginErr != nil {
		fmt.Fprintln(w, loginErr)
	} else {
		fmt.Fprintln(w, loggedUser)
	}
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	user := new(SignupForm)
	if err := parse(r, user); err != nil {
		panic(err)
	}

	createUser := models.User{Name: user.Name, Email: user.Email, Password: user.Password}
	u.UserService.Create(&createUser)
	fmt.Println("email and password", user.Email, user.Password)
	u.NewView.Render(w, nil)
}
