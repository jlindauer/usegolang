package controllers

import (
  "fmt"
  "net/http"
  "strconv"
  "github.com/gorilla/mux"
  "github.com/jlindauer/usegolang/models"
  "github.com/jlindauer/usegolang/views"
  "github.com/jlindauer/usegolang/context"
)

const (
  IndexGalleries = "index_galleries"
  ShowGallery    = "show_gallery"
  EditGallery    = "edit_gallery"
  maxMultipartMem = 1 << 20
)

type Galleries struct {
  New       *views.View
  ShowView  *views.View
  EditView  *views.View
  IndexView *views.View
  gs        models.GalleryService
  is        models.ImageService
  r         *mux.Router
}

type GalleryForm struct {
  Title string `schema:"title"`
}

func NewGalleries(gs models.GalleryService, is models.ImageService, r *mux.Router) *Galleries {
  return &Galleries{
    New:       views.NewView("bootstrap", "galleries/new"),
    ShowView:  views.NewView("bootstrap", "galleries/show"),
    EditView:  views.NewView("bootstrap", "galleries/edit"),
    IndexView: views.NewView("bootstrap", "galleries/index"),
    gs:        gs,
    is:        is,
    r:         r,
  }
}

// POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
  var vd views.Data
  var form GalleryForm
  if err := parseForm(r, &form); err != nil {
    vd.SetAlert(err)
    g.New.Render(w, r, vd)
    return
  }
  user := context.User(r.Context())

  gallery := models.Gallery{
    Title:  form.Title,
    UserID: user.ID,
  }
  if err := g.gs.Create(&gallery); err != nil {
    vd.SetAlert(err)
    g.New.Render(w, r, vd)
    return
  }

  url, err := g.r.Get(EditGallery).URL("id",
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

// POST /galleries/:id/images
func (g *Galleries) ImageUpload(w http.ResponseWriter, r *http.Request) {
  gallery, err := g.galleryByID(w, r)
  if err != nil {
    return
  }
  user := context.User(r.Context())
  if gallery.UserID != user.ID {
    http.Error(w, "Gallery not found", http.StatusNotFound)
    return
  }
  var vd views.Data
  vd.Yield = gallery
  err = r.ParseMultipartForm(maxMultipartMem)
  if err != nil {
    // If we can't parse the form just render an error alert
    // on the edit gallery page.
    vd.SetAlert(err)
    g.EditView.Render(w, r, vd)
    return
  }
  // Iterate over uploaded files to process them.
  files := r.MultipartForm.File["images"]
  for _, f := range files {
    // Open the uploaded file
    file, err := f.Open()
    if err != nil {
      vd.SetAlert(err)
      g.EditView.Render(w, r, vd)
      return
    }
    defer file.Close()
    // Create the image
    err = g.is.Create(gallery.ID, file, f.Filename)
    if err != nil {
      vd.SetAlert(err)
      g.EditView.Render(w, r, vd)
      return
    }
  }
  // Success alert
  vd.Alert = &views.Alert{
    Level:   views.AlertLvlSuccess,
    Message: "Images successfully uploaded!",
  }
  g.EditView.Render(w, r, vd)
}

// POST /galleries/:id/images/:filename/delete
func (g *Galleries) ImageDelete(w http.ResponseWriter, r *http.Request) {
  gallery, err := g.galleryByID(w, r)
  if err != nil {
    return
  }
  user := context.User(r.Context())
  if gallery.UserID != user.ID {
    http.Error(w, "You do not have permission to edit this gallery or image", http.StatusForbidden)
    return
  }
  // Get the filename from the url
  filename := mux.Vars(r)["filename"]
  // Build the Image model
  i := models.Image{
    Filename:  filename,
    GalleryID: gallery.ID,
  }
  // Try to delete the image
  err = g.is.Delete(&i)
  if err != nil {
    // Render the edit page with any errors
    var vd views.Data
    vd.Yield = gallery
    vd.SetAlert(err)
    g.EditView.Render(w, r, vd)
    return
  }
  // If successful redirect to the edit gallery page
  url, err := g.r.Get(EditGallery).URL("id", fmt.Sprintf("%v", gallery.ID))
  if err != nil {
    http.Redirect(w, r, "/galleries", http.StatusFound)
    return
  }
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
  g.ShowView.Render(w, r, vd)
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
  g.EditView.Render(w, r, vd)
}

// GET /galleries
func (g *Galleries) Index(w http.ResponseWriter, r *http.Request) {
  user := context.User(r.Context())
  galleries, err := g.gs.ByUserID(user.ID)
  if err != nil {
    http.Error(w, "Something went wrong.", http.StatusInternalServerError)
    return
  }
  var vd views.Data
  vd.Yield = galleries
  g.IndexView.Render(w, r, vd)
}

// POST /galleries/:id/update
func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {
  gallery, err := g.galleryByID(w, r)
  if err != nil {
    return
  }
  user := context.User(r.Context())
  if gallery.UserID != user.ID {
    http.Error(w, "Gallery not found", http.StatusNotFound)
    return
  }
  var vd views.Data
  vd.Yield = gallery
  var form GalleryForm
  if err := parseForm(r, &form); err != nil {
    // If there is an error we are going to render
    // the EditView again with an alert message
    vd.SetAlert(err)
    g.EditView.Render(w, r, vd)
    return
  }
  gallery.Title = form.Title
  err = g.gs.Update(gallery)
  if err != nil {
    vd.SetAlert(err)
  }else{
    vd.Alert = &views.Alert{
      Level:   views.AlertLvlSuccess,
      Message: "Gallery updated successfully!",
    }
  }
  g.EditView.Render(w, r, vd)
}

// POST /galleries/:id/delete
func (g *Galleries) Delete(w http.ResponseWriter, r *http.Request) {
  // Lookup the gallery using the galleryByID we wrote ealier
  gallery, err := g.galleryByID(w, r)
  if err != nil {
    return
  }
  user := context.User(r.Context())
  if gallery.UserID != user.ID {
    http.Error(w, "You do not have permission to edit this gallery", http.StatusForbidden)
    return
  }

  var vd views.Data
  err = g.gs.Delete(gallery.ID)
  if err != nil {
    vd.SetAlert(err)
    vd.Yield = gallery
    g.EditView.Render(w, r, vd)
    return
  }
  url, err := g.r.Get(IndexGalleries).URL()
  if err != nil {
    http.Redirect(w, r, "/", http.StatusFound)
    return
  }
  http.Redirect(w, r, url.Path, http.StatusFound)
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
  images, _ := g.is.ByGalleryID(gallery.ID)
  gallery.Images = images
  return gallery, nil
}
