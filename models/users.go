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
	db   *gorm.DB
	hmac hash.HMAC
}

type userValidator struct {
	UserDB
}

var _ UserService = &userService{}
var _ UserDB = &userValidator{}

func NewUserService(connectionInfo string) (*userService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(hmacSecretKey)

	ug := userGorm{db: db, hmac: hmac}

	return &userService{UserDB: userValidator{UserDB: &ug}}, nil
}

func (ug *userGorm) Update(user *User) error {
	err := ug.db.Save(user).Error
	if err != nil {
		return err
	}
	if user.Remember != "" {
		user.RememberHash = ug.hmac.Hash(user.Remember)
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
	if id == 0 {
		return ErrInvalidId
	}
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

const passwordPepper = "cC242xTzSG!6j!mWd2N3Vg3!!Q38wunu23a6YBUTm@e**GyP@!CyAzjW7JcR7*p!^524sNxs9H7RQkh3^xH3Q4eSFQtQNqnXqW!"
const hmacSecretKey = "cC242xTzSG!6j!mWd2N3Vh4!!Q38wunu23a6YBUTm@e**GyP@!CyAzjW7JcR7*p!^524sNxs9H7RQkh3^xH3Q4eSFQtQNqnXqW!"

// Create provider user
func (ug *userGorm) Create(user *User) error {
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
	user.RememberHash = ug.hmac.Hash(user.Remember)
	return ug.db.Create(user).Error
}

func (ug *userGorm) ByRemember(token string) (*User, error) {
	var user User
	hashedToken := ug.hmac.Hash(token)
	err := first(ug.db.Where("remember_hash = ?", hashedToken), &user)
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
