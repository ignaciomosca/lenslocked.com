package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
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
	fmt.Println("files", files)
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

func (v *View) Render(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	switch data.(type) {
	case Data:
		//do nothing
	default:
		data = Data{
			Yield: data,
		}

	}
	var buf bytes.Buffer
	if err := v.Template.ExecuteTemplate(&buf, v.Layout, data); err != nil {
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
	v.Render(w, nil)
}

type View struct {
	Template *template.Template
	Layout   string
}
