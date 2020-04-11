package models

import "strings"

var (
	ErrNotFound             modelError = "models: record not found"
	ErrInvalidID            modelError = "models: Id must be greater than 0"
	InvalidPassword         modelError = "models: Password is invalid"
	EmptyPassword           modelError = "models: Password is empty"
	EmptyEmail              modelError = "models: Email is empty"
	InvalidEmail            modelError = "models: Email is invalid"
	EmailAlreadyTaken       modelError = "Email address is already taken"
	PasswordTooShort        modelError = "Password must be at least 8 characters long"
	ErrRememberHashTooShort modelError = "Remember hash too short"
	InvalidHash             modelError = "Remember hash is invalid"
	ErrTokenInvalid         modelError = "models: Reset password token is invalid"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}
