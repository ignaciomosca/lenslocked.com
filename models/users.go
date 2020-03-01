package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialect/postgres"
)

var (
	ErrNotFound  = errors.New("models: Resource not found")
	ErrInvalidId = errors.New("models: Id must be greater than 0")
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}

	return &UserService{db: db}, nil
}

func (us *UserService) Update(user *User) error {
	err := us.db.Save(user).Error
	if err != nil {
		return err
	}
	return nil
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
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

// Create provider user
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

type User struct {
	gorm.Model
	Name  string `gorm:"not null;unique_index`
	Email string
}