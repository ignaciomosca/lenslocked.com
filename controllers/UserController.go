package controllers

import (
	"fmt"
	"net/http"
	"time"

	"lenslocked.com/context"
	"lenslocked.com/email"
	"lenslocked.com/models"
	"lenslocked.com/rand"
	"lenslocked.com/views"
)

// NewUser creates the controller for the operations linked with user management
func NewUser(us models.UserService, emailer *email.Client) Users {
	return Users{
		NewView:      views.NewFiles("bootstrap", "users/new"),
		LoginView:    views.NewFiles("bootstrap", "users/login"),
		ForgotPwView: views.NewFiles("bootstrap", "users/forgot_pw"),
		ResetPwView:  views.NewFiles("bootstrap", "users/reset_pw"),
		UserService:  us,
		emailer:      emailer,
	}
}

type Users struct {
	NewView      *views.View
	LoginView    *views.View
	ForgotPwView *views.View
	ResetPwView  *views.View
	UserService  models.UserService
	emailer      *email.Client
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
	var form SignupForm
	parseURLParams(r, &form)
	u.NewView.Render(w, r, form)
}

func (u *Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	if err := r.ParseForm(); err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}
	login := new(LoginForm)
	if err := parse(r, login); err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}
	loggedUser, loginErr := u.UserService.Login(login.Email, login.Password)
	if loginErr != nil {
		vd.SetAlert(loginErr)
		u.LoginView.Render(w, r, vd)
		return
	}
	u.signIn(w, loggedUser)
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (u *Users) SignOut(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	if user != nil {
		cookie := http.Cookie{Name: "remember_token", Value: "", Expires: time.Now(), HttpOnly: true}
		http.SetCookie(w, &cookie)
	}
	token, _ := rand.RememberToken()
	user.Remember = token
	u.UserService.Update(user)
	http.Redirect(w, r, "/", http.StatusFound)
}

// Create is used to create a new user account
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form SignupForm
	vd.Yield = &form
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
	u.emailer.Welcome(user.Name, user.Email)
	err := u.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

// ResetPwForm is used to process the forgot password form
// and the reset password form.
type ResetPwForm struct {
	Email    string `schema:"email"`
	Token    string `schema:"token"`
	Password string `schema:"password"`
}

// POST /forgot
func (u *Users) InitiateReset(w http.ResponseWriter, r *http.Request) {
	// TODO: Process the forgot password form and iniiate that process
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parse(r, &form); err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}

	token, err := u.UserService.InitiateReset(form.Email)
	if err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}

	err = u.emailer.ResetPw(form.Email, token)
	if err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}

	views.RedirectAlert(w, r, "/reset", http.StatusFound, views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Instructions for resetting your password have been emailed to you.",
	})
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

// ResetPw displays the reset password form and has a method
// so that we can prefill the form data with a token provided
// via the URL query params.
//
// GET /reset
func (u *Users) ResetPw(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parseURLParams(r, &form); err != nil {
		vd.SetAlert(err)
	}
	u.ResetPwView.Render(w, r, vd)
}

// CompleteReset processed the reset password form
//
// POST /reset
func (u *Users) CompleteReset(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parse(r, &form); err != nil {
		vd.SetAlert(err)
		u.ResetPwView.Render(w, r, vd)
		return
	}

	user, err := u.UserService.CompleteReset(form.Token, form.Password)
	if err != nil {
		vd.SetAlert(err)
		u.ResetPwView.Render(w, r, vd)
		return
	}

	u.signIn(w, user)
	views.RedirectAlert(w, r, "/galleries", http.StatusFound, views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Your password has been reset and you have been logged in!",
	})
}
