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

const (
	maxMemory = 52428800 // 50 MB (50 << 20)
)

// NewGallery creates a new GalleryService with it's associated views
func NewGallery(gs models.GalleryService, is models.ImageService, r *mux.Router) Galleries {
	return Galleries{
		NewView:        views.NewFiles("bootstrap", "galleries/new"),
		ShowView:       views.NewFiles("bootstrap", "galleries/show"),
		EditView:       views.NewFiles("bootstrap", "galleries/edit"),
		IndexView:      views.NewFiles("bootstrap", "galleries/index"),
		ImageService:   is,
		GalleryService: gs,
		router:         r,
	}
}

// New renders the view to create a new Gallery
func (g *Galleries) New(w http.ResponseWriter, r *http.Request) {
	g.NewView.Render(w, r, nil)
}

// Galleries contains views and services needed to interact with the galleries table
type Galleries struct {
	NewView        *views.View
	ShowView       *views.View
	EditView       *views.View
	IndexView      *views.View
	GalleryService models.GalleryService
	ImageService   models.ImageService
	router         *mux.Router
}

type GalleryForm struct {
	Title string `schema:"title"`
}

// GET /galleries
func (g *Galleries) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	galleries, err := g.GalleryService.ByUserId(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	var vd views.Data
	vd.Yield = galleries
	g.IndexView.Render(w, r, vd)
}

// GET /galleries/:id
func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.getGalleryByID(w, r)
	if err != nil {
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(w, r, vd)
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
	g.EditView.Render(w, r, vd)
}

// POST /galleries/:id/images
func (g *Galleries) ImageUpload(w http.ResponseWriter, r *http.Request) {
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
	err = r.ParseMultipartForm(maxMemory)
	if err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}

	files := r.MultipartForm.File["images"]
	for _, f := range files {
		file, err := f.Open()
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(w, r, vd)
			return
		}
		defer file.Close()
		err = g.ImageService.Create(gallery.ID, file, f.Filename)
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(w, r, vd)
			return
		}
	}
	url, err := g.router.Get("edit_gallery").URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		http.Redirect(w, r, "/galleries", http.StatusFound)
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
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
		g.EditView.Render(w, r, vd)
		return
	}
	gallery.Title = form.Title
	err = g.GalleryService.Update(gallery)
	if err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}
	vd.Alert = &views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Gallery succesfully updated",
	}
	g.EditView.Render(w, r, vd)
}

// GET /galleries/:id/delete
func (g *Galleries) Delete(w http.ResponseWriter, r *http.Request) {
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
	err = g.GalleryService.Delete(gallery.ID)
	if err != nil {
		vd.SetAlert(err)
		vd.Yield = gallery
		g.EditView.Render(w, r, vd)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (g *Galleries) getGalleryByID(w http.ResponseWriter, r *http.Request) (*models.Gallery, error) {
	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return nil, err
	}
	gallery, err := g.GalleryService.ById(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	images, _ := g.ImageService.ByGalleryID(uint(id))
	gallery.Images = images
	return gallery, nil
}

// Create is used to create a new gallery
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parse(r, &form); err != nil {
		vd.SetAlert(err)
		g.NewView.Render(w, r, vd)
		return
	}

	user := context.User(r.Context())
	log.Println("Create got the user", user)
	if user == nil {
		log.Println("User not found")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	gallery := models.Gallery{Title: form.Title, UserID: user.ID}
	if err := g.GalleryService.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.NewView.Render(w, r, vd)
		return
	}
	url, err := g.router.Get("show_gallery").URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/galleries", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)

}
