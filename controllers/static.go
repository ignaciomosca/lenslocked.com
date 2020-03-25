package controllers

import "lenslocked.com/views"

// NewStatic creates a struct with the static pages
func NewStatic() *Static {
	return &Static{
		HomeView:     views.NewFiles("bootstrap", "static/home"),
		ContactView:  views.NewFiles("bootstrap", "static/contact"),
		FAQView:      views.NewFiles("bootstrap", "static/faq"),
		NotFoundView: views.NewFiles("bootstrap", "static/notfound"),
	}
}

// Static contains the views for the static pages
type Static struct {
	HomeView     *views.View
	ContactView  *views.View
	FAQView      *views.View
	NotFoundView *views.View
}
