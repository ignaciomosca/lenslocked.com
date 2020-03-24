package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/models"
	"lenslocked.com/rand"
	"lenslocked.com/views"
)

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
		fmt.Println("Login Finished. User", loggedUser)
		u.signIn(w, loggedUser)
		http.Redirect(w, r, "/cookietest", http.StatusFound)
	}
}

// Create is used to create a new user account
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form SignupForm
	if err := parse(r, &form); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, vd)
		return
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.UserService.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, vd)
		return
	}
	err := u.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
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
			fmt.Println("Error 2", err)
			return err
		}
	}

	fmt.Println("Cookie", user.Remember)
	cookie := http.Cookie{Name: "remember_token", Value: user.Remember, HttpOnly: true}
	http.SetCookie(w, &cookie)
	return nil
}
