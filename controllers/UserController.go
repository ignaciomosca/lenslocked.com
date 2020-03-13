package controllers

import (
	"fmt"
	"log"
	"net/http"

	"lenslocked.com/models"
	"lenslocked.com/rand"
	"lenslocked.com/views"
)

func NewUser(connectionInfo string) Users {
	us, err := models.NewUserService(connectionInfo)
	if err != nil {
		panic(err)
	}

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
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
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
		u.signIn(w, loggedUser)
		http.Redirect(w, r, "/cookieTest", http.StatusFound)
	}
}

// Create is used to create a new user account
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	var form SignupForm
	if err := parse(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		u.NewView.Render(w, vd)
		return
	}

	createUser := models.User{Name: form.Name, Email: form.Email, Password: form.Password}
	if err := u.UserService.Create(&createUser); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		u.NewView.Render(w, vd)
		return
	}

	if err := u.signIn(w, &createUser); err != nil {
		log.Println(err)
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlDanger,
			Message: err.Error(),
		}
		u.NewView.Render(w, vd)
		return
	}
	err := u.signIn(w, &createUser)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/cookieTest", http.StatusFound)
}

func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		panic(err)
	}
	user, err := u.UserService.ByRemember(cookie.Value)
	if err != nil {
		return
	}
	fmt.Fprintln(w, user)
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
			return err
		}
	}

	cookie := http.Cookie{Name: "remember_token", Value: user.Remember, HttpOnly: true}
	http.SetCookie(w, &cookie)
	fmt.Fprintln(w, u)
	return nil
}
