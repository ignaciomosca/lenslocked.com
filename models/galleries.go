package models

import (
	"github.com/jinzhu/gorm"
)

// Gallery is our image container resources that visitors view
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}

type GalleryService interface {
	GalleryDB
}
type GalleryDB interface {
	ById(id uint) (*Gallery, error)
	Create(gallery *Gallery) error
}

// ById lookups a gallery by id.
// If no gallery exists returns an error
func (gg *galleryGorm) ById(id uint) (*Gallery, error) {
	var gallery Gallery
	err := gg.db.Where("id = ?", id).First(&gallery).Error
	switch err {
	case nil:
		return &gallery, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

type galleryValidator struct {
	GalleryDB
}

func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: galleryValidator{
			&galleryGorm{db},
		},
	}
}

type galleryService struct {
	GalleryDB
}

type galleryGorm struct {
	db *gorm.DB
}

var _ GalleryDB = &galleryGorm{}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	gg.db.AutoMigrate(&Gallery{})
	return gg.db.Create(gallery).Error
}
