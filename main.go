package main

import (
  "net/http"
  "github.com/gorilla/mux"
  "github.com/jlindauer/usegolang/views"
)

var homeView *views.View
var contactView *views.View

func home(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html")
  must(homeView.Render(w, nil))
}

func contact(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html")
  must(homeView.Render(w, nil))
}

func must(err error) {
  if err != nil {
    panic(err)
  }
}

func main() {
  homeView = views.NewView("bootstrap", "views/home.gohtml")
  contactView = views.NewView("bootstrap", "views/contact.gohtml")

  r := mux.NewRouter()
  r.HandleFunc("/", home)
  r.HandleFunc("/contact", contact)
  http.ListenAndServe(":3000", r)
}
