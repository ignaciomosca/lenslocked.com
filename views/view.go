package views

import (
	"fmt"
	"html/template"
	"path/filepath"
)

var (
	LayoutDir         string = "views/layouts/"
	TemplateExtension string = ".gohtml"
)

func NewFiles(layout string, files ...string) *View {
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

func fetchFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExtension)
	if err != nil {
		panic(err)
	}
	return files
}

type View struct {
	Template *template.Template
	Layout   string
}
