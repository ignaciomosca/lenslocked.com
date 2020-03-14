package models

import (
	"github.com/jinzhu/gorm"
)

func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &Services{
		User: NewUserService(db),
		db:   db,
	}, nil
}

type Services struct {
	User    UserService
	Gallery GalleryService
	db      *gorm.DB
}

// Closes database connection
func (s *Services) Close() error {
	return s.db.Close()
}
