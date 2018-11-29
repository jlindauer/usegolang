package controllers

import (
  "fmt"
  "net/http"
  "github.com/jlindauer/usegolang/models"
  "github.com/jlindauer/usegolang/views"
)

type Galleries struct {
  New *views.View
  gs  models.GalleryService
}

type GalleryForm struct {
  Title string `schema:"title"`
}

func NewGalleries(gs models.GalleryService) *Galleries {
  return &Galleries{
    New: views.NewView("bootstrap", "galleries/new"),
    gs:  gs,
  }
}

// POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
  var vd views.Data
  var form GalleryForm
  if err := parseForm(r, &form); err != nil {
    vd.SetAlert(err)
    g.New.Render(w, vd)
    return
  }
  gallery := models.Gallery{
    Title: form.Title,
  }
  if err := g.gs.Create(&gallery); err != nil {
    vd.SetAlert(err)
    g.New.Render(w, vd)
    return
  }
  fmt.Fprintln(w, gallery)
}
