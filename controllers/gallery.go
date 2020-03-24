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

func NewGallery(gs models.GalleryService, r *mux.Router) Galleries {
	return Galleries{
		NewView:        views.NewFiles("bootstrap", "galleries/new"),
		ShowView:       views.NewFiles("bootstrap", "galleries/show"),
		GalleryService: gs,
		router:         r,
	}
}

func (g *Galleries) New(w http.ResponseWriter, r *http.Request) {
	g.NewView.Render(w, nil)
}

type Galleries struct {
	NewView        *views.View
	ShowView       *views.View
	GalleryService models.GalleryService
	router         *mux.Router
}

type GalleryForm struct {
	Title string `schema:"title"`
}

// GET /galleries/:id
func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	res, err := g.GalleryService.ById(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	g.ShowView.Render(w, res)
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
