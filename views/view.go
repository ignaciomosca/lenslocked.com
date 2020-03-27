package views

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"lenslocked.com/context"
)

var (
	LayoutDir         string = "views/layouts/"
	TemplateDir       string = "views/"
	TemplateExtension string = ".gohtml"
)

func NewFiles(layout string, files ...string) *View {
	addTemplateFiles(files)
	addTemplateExt(files)
	files = append(files, fetchFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout}
}

func addTemplateFiles(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExtension
	}
}

func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
	default:
		vd = Data{
			Yield: data,
		}

	}
	vd.User = context.User(r.Context())

	var buf bytes.Buffer
	if err := v.Template.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
		log.Println("error", err)
		http.Error(w, "Something went wrong rendering a page", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func fetchFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExtension)
	if err != nil {
		panic(err)
	}
	return files
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	v.Render(w, r, nil)
}

type View struct {
	Template *template.Template
	Layout   string
}
