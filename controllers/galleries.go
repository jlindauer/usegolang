package controllers

import (
  "net/http"
  "strconv"
  "github.com/gorilla/mux"
  "github.com/jlindauer/usegolang/models"
  "github.com/jlindauer/usegolang/views"
  "github.com/jlindauer/usegolang/context"
)

const (
  ShowGallery = "show_gallery"
)

type Galleries struct {
  New      *views.View
  ShowView *views.View
  EditView *views.View
  gs       models.GalleryService
  r        *mux.Router
}

type GalleryForm struct {
  Title string `schema:"title"`
}

func NewGalleries(gs models.GalleryService, r *mux.Router) *Galleries {
  return &Galleries{
    New:      views.NewView("bootstrap", "galleries/new"),
    ShowView: views.NewView("bootstrap", "galleries/show"),
    EditView: views.NewView("bootstrap", "galleries/edit"),
    gs:       gs,
    r:        r,
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
  user := context.User(r.Context())

  gallery := models.Gallery{
    Title:  form.Title,
    UserID: user.ID,
  }
  if err := g.gs.Create(&gallery); err != nil {
    vd.SetAlert(err)
    g.New.Render(w, vd)
    return
  }

  url, err := g.r.Get(ShowGallery).URL("id",
    strconv.Itoa(int(gallery.ID)))
  if err != nil {
    http.Redirect(w, r, "/", http.StatusFound)
    return
  }
  // If no errors, use the URL we just created and redirect
  // to the path portion of that URL. We don't need the
  // entire URL because your application might be hosted at
  // localhost:3000, or it might be at lenslocked.com. By
  // only using the path our code is agnostic to that detail.
  http.Redirect(w, r, url.Path, http.StatusFound)
}

// GET /galleries/:id
func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
  gallery, err := g.galleryByID(w, r)
  if err != nil {
    return
  }
  var vd views.Data
  vd.Yield = gallery
  g.ShowView.Render(w, vd)
}

// GET /galleries/:id/edit
func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
  gallery, err := g.galleryByID(w, r)
  if err != nil {
    return
  }
  // A user needs to be logged in to acces this page, so we can
  // assume that the RequireUser middleware has run and set
  // the user for us in the request context.
  user := context.User(r.Context())
  if gallery.UserID != user.ID {
    http.Error(w, "You do not have permission to edit this gallery", http.StatusForbidden)
    return
  }
  var vd views.Data
  vd.Yield = gallery
  g.EditView.Render(w, vd)
}

func (g *Galleries) galleryByID(w http.ResponseWriter, r *http.Request) (*models.Gallery, error) {
  vars := mux.Vars(r)
  idStr := vars["id"]
  id, err := strconv.Atoi(idStr)
  if err != nil {
    http.Error(w, "Invalid gallery ID", http.StatusNotFound)
    return nil, err
  }
  gallery, err := g.gs.ByID(uint(id))
  if err != nil {
    switch err {
    case models.ErrNotFound:
      http.Error(w, "Gallery not found", http.StatusNotFound)
    default:
      http.Error(w, "Whoops! Something went wrong", http.StatusInternalServerError)
    }
    return nil, err
  }
  return gallery, nil
}
