package controllers

import "lenslocked.com/views"

func NewStatic() *Static {
	return &Static{
		HomeView:     views.NewFiles("bootstrap", "static/home"),
		ContactView:  views.NewFiles("bootstrap", "static/contact"),
		FAQView:      views.NewFiles("bootstrap", "static/faq"),
		NotFoundView: views.NewFiles("bootstrap", "static/notfound"),
	}
}

type Static struct {
	HomeView     *views.View
	ContactView  *views.View
	FAQView      *views.View
	NotFoundView *views.View
}
