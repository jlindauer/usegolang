package controllers

import (
  "github.com/jlindauer/usegolang/models"
  "github.com/jlindauer/usegolang/views"
)

type Galleries struct {
  New *views.View
  gs  models.GalleryService
}

func NewGalleries(gs models.GalleryService) *Galleries {
  return &Galleries{
    New: views.NewView("bootstrap", "galleries/new"),
    gs:  gs,
  }
}
