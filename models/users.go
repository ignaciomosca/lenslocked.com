package models

import (
	"errors"

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
)

type UserService struct {
	db   *gorm.DB
	hmac hash.HMAC
}

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(hmacSecretKey)
	return &UserService{db: db, hmac: hmac}, nil
}

func (us *UserService) Update(user *User) error {
	err := us.db.Save(user).Error
	if err != nil {
		return err
	}
	if user.Remember != "" {
		user.RememberHash = us.hmac.Hash(user.Remember)
	}
	return us.db.Save(user).Error
}

func (us *UserService) Login(email, password string) (*User, error) {
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

func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	err := us.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidId
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

func (us *UserService) Close() error {
	return us.db.Close()
}

// ById lookups a user by id.
// If no user exists returns an error
func (us *UserService) ById(id uint) (*User, error) {
	var user User
	err := us.db.Where("id = ?", id).First(&user).Error
	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

const passwordPepper = "cC242xTzSG!6j!mWd2N3Vg3!!Q38wunu23a6YBUTm@e**GyP@!CyAzjW7JcR7*p!^524sNxs9H7RQkh3^xH3Q4eSFQtQNqnXqW!"
const hmacSecretKey = "cC242xTzSG!6j!mWd2N3Vh4!!Q38wunu23a6YBUTm@e**GyP@!CyAzjW7JcR7*p!^524sNxs9H7RQkh3^xH3Q4eSFQtQNqnXqW!"

// Create provider user
func (us *UserService) Create(user *User) error {
	passwordBytes := []byte(user.Password + passwordPepper)
	hashPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashPassword)
	user.Password = ""
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	user.RememberHash = us.hmac.Hash(user.Remember)
	return us.db.Create(user).Error
}

func (us *UserService) ByRemember(token string) (*User, error) {
	var user User
	hashedToken := us.hmac.Hash(token)
	err := first(us.db.Where("remember_hash = ?", hashedToken), &user)
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
