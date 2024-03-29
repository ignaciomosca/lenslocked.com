package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const RememberTokenBytes = 32

// RememberToken is a helper function designed to generate
// remember tokens of a predetermined byte size.
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}

// String will generate a byte slice of size nBytes and then
// return a string that is the base64 URL encoded version
// of that byte slice
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Bytes will help us generate n random bytes, or will
// return an error if there was one. This uses the
// crypto/rand package so it is safe to use with things
// like remember tokens.
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Bytes", err)
		return nil, err
	}
	return b, nil
}

// NBytes base64-decode a string and counts the number of bytes
func NBytes(base64String string) (int, error) {
	b, err := base64.URLEncoding.DecodeString(base64String)
	if err != nil {
		fmt.Println("NBytes", err)
		return -1, err
	}
	return len(b), nil
}
