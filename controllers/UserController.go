package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/models"
	"lenslocked.com/rand"
	"lenslocked.com/views"
)

// NewUser creates the controller for the operations linked with user management
func NewUser(us models.UserService) Users {
	return Users{
		NewView:     views.NewFiles("bootstrap", "users/new"),
		LoginView:   views.NewFiles("bootstrap", "users/login"),
		UserService: us,
	}
}

type Users struct {
	NewView     *views.View
	LoginView   *views.View
	UserService models.UserService
}

// SignupForm represents the data in the signup form
type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// LoginForm represents the data in the login form
type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, r, nil)
}

func (u *Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	if err := r.ParseForm(); err != nil {
		vd.SetAlert(err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	login := new(LoginForm)
	if err := parse(r, login); err != nil {
		vd.SetAlert(err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	loggedUser, loginErr := u.UserService.Login(login.Email, login.Password)
	if loginErr != nil {
		vd.SetAlert(loginErr)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	u.signIn(w, loggedUser)
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

// Create is used to create a new user account
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form SignupForm
	if err := parse(r, &form); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.UserService.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}
	err := u.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.UserService.Update(user)
		if err != nil {
			fmt.Println("Error 2", err)
			return err
		}
	}

	fmt.Println("Cookie", user.Remember)
	cookie := http.Cookie{Name: "remember_token", Value: user.Remember, HttpOnly: true}
	http.SetCookie(w, &cookie)
	return nil
}
