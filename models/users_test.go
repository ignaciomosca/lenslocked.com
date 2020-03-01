package models

import (
	"fmt"
	"testing"
)

func createUserService() (*UserService, error) {
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "postgres"
		dbName   = "postgres"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbName=%s sslMode=disable", host, port, user, password, dbName)
	us, err := NewUserService(psqlInfo)
	if err != nil {
		return nil, err
	}
	us.db.LogMode(false)
	return us, nil

}

func testUserService(t *testing.T) {
	us, err := createUserService()
	if err != nil {
		t.Fatal(err)
	}
	user := User{Name: "Michael Scott", Email: "michaelscott@dundermiffin.com"}
	err = us.Create(&user)
	if err != nil {
		t.Fatal(err)
	}
	if user.ID == 0 {
		t.Errorf("Expected user id > 0", user.ID)
	}

}
