package models

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
	"lenslocked.com/hash"
	"lenslocked.com/rand"
)

type userService struct {
	UserDB
	passwordPepper string
}

// UserDB interacts with the user database
type UserDB interface {
	ById(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	error
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
	hmac           hash.HMAC
	emailRegex     *regexp.Regexp
	passwordPepper string
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

var _ UserDB = &userGorm{}

// Create provider user
func (uv *userValidator) Create(user *User) error {
	err := runUserValFuncs(user,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.defaultRemember,
		// 	uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emaillFormat,
		uv.emailIsAvailable)
	if err != nil {
		fmt.Println("error", err)
		return err
	}
	return uv.UserDB.Create(user)
}

func (uv *userValidator) Update(user *User) error {
	err := runUserValFuncs(user,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emaillFormat,
		uv.emailIsAvailable)
	if err != nil {
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

func (uv *userValidator) emailIsAvailable(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	if err != nil && err == ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	if user.ID != existing.ID {
		return EmailAlreadyTaken
	}
	return nil

}

func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return PasswordTooShort
	}
	return nil
}

func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return EmptyPassword
	}
	return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return InvalidPassword
	}
	return nil
}

func (uv *userValidator) rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return InvalidHash
	}
	return nil
}

func (uv *userValidator) rememberMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}
	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return err
	}
	if n < 32 {
		return ErrRememberHashTooShort
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

func newUserValidator(udb UserDB, hmac hash.HMAC, pepper string) *userValidator {
	return &userValidator{
		UserDB:         udb,
		hmac:           hmac,
		emailRegex:     regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
		passwordPepper: pepper,
	}
}

func NewUserService(db *gorm.DB, pepper, hmacKey string) *userService {
	ug := &userGorm{db: db}
	hmac := hash.NewHMAC(hmacKey)
	uv := newUserValidator(ug, hmac, pepper)

	return &userService{
		UserDB:         uv,
		passwordPepper: pepper,
	}
}

func (ug *userGorm) Update(user *User) error {
	err := ug.db.Save(user).Error
	if err != nil {
		fmt.Println("Updating")
		return err
	}
	return ug.db.Save(user).Error
}

func (us *userService) Login(email, password string) (*User, error) {
	user, err := us.ByEmail(email)
	if err != nil {
		return nil, ErrNotFound
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password+us.passwordPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			fmt.Println("Failed login")
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
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
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
	passwordBytes := []byte(user.Password + uv.passwordPepper)
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
		return ErrInvalidID
	}
	return nil
}

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

func (ug *userGorm) Error() string {
	return "Something went wrong with GORM"
}

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
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
