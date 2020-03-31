package models

import (
	"fmt"
	"io"
	"os"

	"github.com/jinzhu/gorm"
)

type ImageService interface {
	Create(galleryId uint, r io.ReadCloser, filename string) error
	//ByGalleryID(galleryID uint) []string
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

func (is *imageService) mkImagePath(galleryID uint) (string, error) {
	// Create the directory to contain our images
	galleryPath := fmt.Sprintf("images/galleries/%v/", galleryID)
	err := os.MkdirAll(galleryPath, 0755)
	if err != nil {
		return "", err
	}
	return galleryPath, nil
}
