package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jinzhu/gorm"
)

type Image struct {
	GalleryID uint
	Filename  string
}

func (i *Image) Path() string {
	return "/" + i.RelativePath()
}

func (i *Image) RelativePath() string {
	return fmt.Sprintf("images/galleries/%v/%v", i.GalleryID, i.Filename)
}

type ImageService interface {
	Create(galleryId uint, r io.ReadCloser, filename string) error
	ByGalleryID(galleryID uint) ([]Image, error)
	Delete(image *Image) error
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

func (is *imageService) ByGalleryID(galleryID uint) ([]Image, error) {
	path := is.imagePath(galleryID)
	stringNames, err := filepath.Glob(path + "*")
	if err != nil {
		return nil, err
	}
	ret := make([]Image, len(stringNames))
	for i := range stringNames {
		stringNames[i] = strings.Replace(stringNames[i], path, "", 1)
		ret[i] = Image{
			GalleryID: galleryID,
			Filename:  stringNames[i],
		}
	}
	return ret, nil
}

func (is *imageService) Delete(i *Image) error {
	r := os.Remove(i.RelativePath())
	if r != nil {
		panic(r)
		return r
	}
	return r
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
