package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialect/postgres"
)

type User struct {
	gorm.Model
	Name  string `gorm:"not null;unique_index`
	Email string
}
