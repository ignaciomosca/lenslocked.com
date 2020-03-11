package views

import (
	"fmt"
	"html/template"
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

func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	fmt.Println("Data", data)
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

func fetchFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExtension)
	if err != nil {
		panic(err)
	}
	return files
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type Alert struct {
		Level   string
		Message string
	}

	w.Header().Set("Content-Type", "text/html")
	a := Alert{Level: "success", Message: "successfully rendered a dynamic alert"}
	if err := v.Render(w, a); err != nil {
		panic(err)
	}
}

type View struct {
	Template *template.Template
	Layout   string
}
