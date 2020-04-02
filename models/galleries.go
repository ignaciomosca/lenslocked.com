package models

import (
	"log"

	"github.com/jinzhu/gorm"
)

const (
	ErrUserIDRequired   modelError = "models: user ID is required"
	ErrTitleRequired    modelError = "models: title is required"
	ErrGalleryIDInvalid modelError = "models: ID provided was invalid"
)

// Gallery is our image container resources that visitors view
type Gallery struct {
	gorm.Model
	UserID uint    `gorm:"not_null;index"`
	Title  string  `gorm:"not_null"`
	Images []Image `gorm:"-"`
}

func (g *Gallery) ImagesSplitN(n int) [][]Image {
	ret := make([][]Image, n)
	for i := 0; i < n; i++ {
		ret[i] = make([]Image, 0)
	}
	for i, img := range g.Images {
		bucket := i % n
		ret[bucket] = append(ret[bucket], img)
	}
	return ret
}

type GalleryService interface {
	GalleryDB
}
type GalleryDB interface {
	ById(id uint) (*Gallery, error)
	Create(gallery *Gallery) error
	Update(gallery *Gallery) error
	Delete(id uint) error
	ByUserId(userId uint) ([]Gallery, error)
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

func (gv galleryValidator) Create(gallery *Gallery) error {
	err := runGalleryValFns(gallery,
		gv.userIDRequired,
		gv.titleRequired)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Create(gallery)
}

func (gv galleryValidator) Update(gallery *Gallery) error {
	err := runGalleryValFns(gallery,
		gv.userIDRequired,
		gv.titleRequired)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Update(gallery)
}

// Delete will delete the gallery by id
func (gv galleryValidator) Delete(id uint) error {
	var gallery Gallery
	gallery.ID = id
	err := runGalleryValFns(&gallery,
		gv.nonZeroID)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Delete(id)
}

func (gv *galleryValidator) userIDRequired(g *Gallery) error {
	if g.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (gv *galleryValidator) titleRequired(g *Gallery) error {
	if g.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

func (gv *galleryValidator) nonZeroID(gallery *Gallery) error {
	if gallery.ID <= 0 {
		return ErrGalleryIDInvalid
	}
	return nil
}

type galleryService struct {
	GalleryDB
}

type galleryGorm struct {
	db *gorm.DB
}

var _ GalleryDB = &galleryGorm{}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}

func (gg *galleryGorm) Update(gallery *Gallery) error {
	return gg.db.Save(gallery).Error
}

func (gg *galleryGorm) Delete(id uint) error {
	gallery := Gallery{Model: gorm.Model{ID: id}}
	return gg.db.Delete(&gallery).Error
}

func (gg *galleryGorm) ByUserId(userID uint) ([]Gallery, error) {
	var galleries []Gallery
	err := gg.db.Where("user_id = ?", userID).Find(&galleries).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return galleries, nil
}

type galleryValFn func(*Gallery) error

func runGalleryValFns(gallery *Gallery, fns ...galleryValFn) error {
	for _, fn := range fns {
		if err := fn(gallery); err != nil {
			return err
		}
	}
	return nil
}
