package controllers

import (
	"fmt"
	"log"
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

type GalleryForm struct {
	Title string `schema:"title"`
}

// Create is used to create a new gallery
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	var form GalleryForm
	if err := parse(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.NewView.Render(w, vd)
		return
	}

	gallery := models.Gallery{Title: form.Title}
	if err := g.GalleryService.Create(&gallery); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.NewView.Render(w, vd)
		return
	}

	fmt.Fprint(w, gallery)

}
