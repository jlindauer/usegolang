package views

import (
	"html/template"
	"net/http"
	"path/filepath"
	"bytes"
	"io"
	"github.com/jlindauer/usegolang/context"
)

var (
	LayoutDir   string = "views/layouts/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

// Render takes in a ResponseWriter and data, and then executes the template associated
// with the View v, and writes the output to the ResponseWriter
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
	err := v.Template.ExecuteTemplate(&buf, v.Layout, vd)
	if err != nil {
		http.Error(w, "Something went wrong. If the problem " +
			"persists, please email support@usegolang.com",
			http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

// ServeHTTP fulfills the specifications of https://golang.org/pkg/net/http/#ServeMux.ServeHTTP
// It is used to dispatch requests to the handler setup for the router pattern
func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}

// layoutFiles takes in the layout directory file path and Go's template extension and then
// searches with a wildcard to identify all matching files in the directory, the output
// is a slice of strings of the files
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

// addTemplatePath takes in a slice of string representing file paths for templates, and it
// prepends the TemplateDir directory to each string in the slice
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

// addTemplateExt takes in a slice of strings representing file paths for templates
// and it appends the TemplateExt extension to each string in the slice
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}

// NewView takes in the layout filepath as well as all necessary files to create the views
// then the files are parsed and the templates are associated with t. A view is then returned
func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}
