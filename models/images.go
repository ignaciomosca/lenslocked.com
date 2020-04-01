package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"
)

type ImageService interface {
	Create(galleryId uint, r io.ReadCloser, filename string) error
	ByGalleryID(galleryID uint) ([]string, error)
}

type imageService struct {
}

func NewImageService(db *gorm.DB) ImageService {
	return &imageService{}
}

func (is *imageService) Create(galleryId uint, r io.ReadCloser, filename string) error {
	defer r.Close()
	path, err := is.mkImagePath(galleryId)

	dst, err := os.Create(path + filename)
	defer dst.Close()

	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}
	return nil
}

func (is *imageService) ByGalleryID(galleryID uint) ([]string, error) {
	path := is.imagePath(galleryID)
	strings, err := filepath.Glob(path + "*")
	if err != nil {
		return nil, err
	}
	for i := range strings {
		strings[i] = "/" + strings[i]
	}
	return strings, nil
}

func (is *imageService) imagePath(galleryID uint) string {
	return fmt.Sprintf("images/galleries/%v/", galleryID)
}

func (is *imageService) mkImagePath(galleryID uint) (string, error) {
	// Create the directory to contain our images
	galleryPath := is.imagePath(galleryID)
	err := os.MkdirAll(galleryPath, 0755)
	if err != nil {
		return "", err
	}
	return galleryPath, nil
}
