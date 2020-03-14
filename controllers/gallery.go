package controllers

import (
	"net/http"

	"lenslocked.com/models"
	"lenslocked.com/views"
)

func NewGallery(gs models.GalleryService) Galleries {
	return Galleries{
		NewView:        views.NewFiles("bootstrap", "galleries/new"),
		GalleryService: gs,
	}
}

func (g *Galleries) New(w http.ResponseWriter, r *http.Request) {
	g.NewView.Render(w, nil)
}

type Galleries struct {
	NewView        *views.View
	GalleryService models.GalleryService
}
