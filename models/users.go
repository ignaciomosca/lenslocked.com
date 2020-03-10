package models

import (
	"errors"
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
	"lenslocked.com/hash"
	"lenslocked.com/rand"
)

var (
	ErrNotFound     = errors.New("models: Resource not found")
	ErrInvalidId    = errors.New("models: Id must be greater than 0")
	InvalidPassword = errors.New("models: Password is invalid")
	EmptyEmail      = errors.New("models: Email is empty")
	InvalidEmail    = errors.New("models: Email is invalid")
)

type userService struct {
	UserDB
}

// UserDB interacts with the user database
type UserDB interface {
	ById(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	Close() error
}

// UserService is a set of methods to manipulate the user model
type UserService interface {
	// Login Verifies user and password are correct
	// if user/password is correct a user gets returned
	// if not, it will return an error
	Login(email, password string) (*User, error)
	UserDB
}

type userGorm struct {
	db *gorm.DB
}

type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}

// ByRemember will hash the remember token and then call ByRemember on the UserDB layer.
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFuncs(&user, uv.hmacRemember); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

// Create provider user
func (uv *userValidator) Create(user *User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token

	}

	if err := runUserValFuncs(user, uv.bcryptPassword, uv.defaultRemember, uv.hmacRemember, uv.normalizeEmail, uv.requireEmail, uv.emaillFormat); err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

func (uv *userValidator) Update(user *User) error {
	if err := runUserValFuncs(user, uv.bcryptPassword, uv.hmacRemember, uv.normalizeEmail, uv.requireEmail, uv.emaillFormat); err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

// Delete will delete a user with a provided id
func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValFuncs(&user, uv.idGreaterThanZero)
	if err != nil {
		return err
	}

	return uv.UserDB.Delete(id)
}

// ByEmail will normalize the email address before querying the UserDB
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{Email: email}
	err := runUserValFuncs(&user, uv.normalizeEmail)
	if err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(email)
}

func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	if user.Email == "" {
		return EmptyEmail
	}

	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return EmptyEmail
	}
	return nil
}

func (uv *userValidator) emaillFormat(user *User) error {
	if user.Email == "" {
		return nil
	}
	if !uv.emailRegex.MatchString(user.Email) {
		return InvalidEmail
	}
	return nil
}

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

var _ UserService = &userService{}
var _ UserDB = &userValidator{}

func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB:     udb,
		hmac:       hmac,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

func NewUserService(connectionInfo string) (*userService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(hmacSecretKey)

	ug := userGorm{db: db}
	uv := &userValidator{UserDB: &ug, hmac: hmac}

	return &userService{UserDB: uv}, nil
}

func (ug *userGorm) Update(user *User) error {
	err := ug.db.Save(user).Error
	if err != nil {
		return err
	}
	return ug.db.Save(user).Error
}

// func (uv UserValidator) ById(id uint) (*User, error) {
// 	if id <= 0 {
// 		return nil, ErrInvalidId
// 	}
// 	return uv.UserDB.ById(id)
// }

func (us *userService) Login(email, password string) (*User, error) {
	user, err := us.ByEmail(email)
	if err != nil {
		return nil, ErrNotFound
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password+passwordPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, InvalidPassword
		default:
			return nil, err
		}
	} else {
		return user, nil
	}
}

func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	err := ug.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

func (ug *userGorm) Close() error {
	return ug.db.Close()
}

// ById lookups a user by id.
// If no user exists returns an error
func (ug *userGorm) ById(id uint) (*User, error) {
	var user User
	err := ug.db.Where("id = ?", id).First(&user).Error
	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// bcryptPassword will hash a users password with a predefined pepper and bcrypt
// if the password field is not the empty string
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	passwordBytes := []byte(user.Password + passwordPepper)
	hashPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashPassword)
	user.Password = ""
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

func (uv *userValidator) defaultRemember(user *User) error {
	if user.Remember != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

func (uv *userValidator) idGreaterThanZero(user *User) error {
	if user.ID <= 0 {
		return ErrInvalidId
	}
	return nil
}

const passwordPepper = "cC242xTzSG!6j!mWd2N3Vg3!!Q38wunu23a6YBUTm@e**GyP@!CyAzjW7JcR7*p!^524sNxs9H7RQkh3^xH3Q4eSFQtQNqnXqW!"
const hmacSecretKey = "cC242xTzSG!6j!mWd2N3Vh4!!Q38wunu23a6YBUTm@e**GyP@!CyAzjW7JcR7*p!^524sNxs9H7RQkh3^xH3Q4eSFQtQNqnXqW!"

// Create provider user
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

// ByRemember looks a user up with a given rememberHash and returns that user.
// This method expects the remember token to already be hashed.
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

type User struct {
	gorm.Model
	Name         string
	Email        string
	Password     string `gorm:"-"`
	PasswordHash string
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}
