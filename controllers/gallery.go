package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"lenslocked.com/context"
	"lenslocked.com/models"
	"lenslocked.com/views"
)

// NewGallery creates a new GalleryService with it's associated views
func NewGallery(gs models.GalleryService, r *mux.Router) Galleries {
	return Galleries{
		NewView:        views.NewFiles("bootstrap", "galleries/new"),
		ShowView:       views.NewFiles("bootstrap", "galleries/show"),
		EditView:       views.NewFiles("bootstrap", "galleries/edit"),
		GalleryService: gs,
		router:         r,
	}
}

// New renders the view to create a new Gallery
func (g *Galleries) New(w http.ResponseWriter, r *http.Request) {
	g.NewView.Render(w, nil)
}

// Galleries contains views and services needed to interact with the galleries table
type Galleries struct {
	NewView        *views.View
	ShowView       *views.View
	EditView       *views.View
	GalleryService models.GalleryService
	router         *mux.Router
}

type GalleryForm struct {
	Title string `schema:"title"`
}

// GET /galleries/:id
func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.getGalleryByID(w, r)
	if err != nil {
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(w, vd)
}

// GET /galleries/:id/edit
func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.getGalleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.EditView.Render(w, vd)
}

// GET /galleries/:id/update
func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.getGalleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	var form GalleryForm
	if err := parse(r, &form); err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, vd)
		return
	}
	gallery.Title = form.Title
	err = g.GalleryService.Update(gallery)
	if err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, vd)
		return
	}
	vd.Alert = &views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Gallery succesfully updated",
	}
	g.EditView.Render(w, vd)
}

func (g *Galleries) getGalleryByID(w http.ResponseWriter, r *http.Request) (*models.Gallery, error) {
	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return nil, err
	}
	res, err := g.GalleryService.ById(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	return res, nil
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

	user := context.User(r.Context())
	fmt.Println("Create got the user", user)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	gallery := models.Gallery{Title: form.Title, UserID: user.ID}
	if err := g.GalleryService.Create(&gallery); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.NewView.Render(w, vd)
		return
	}
	vd.Yield = gallery
	url, err := g.router.Get("show_gallery").URL(string(gallery.ID))
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)

}