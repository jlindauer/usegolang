package main

import (
  "net/http"
  "github.com/gorilla/mux"
  "github.com/jlindauer/usegolang/views"
  "github.com/jlindauer/usegolang/controllers"
)

var (
  homeView    *views.View
  contactView *views.View
)

func home(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html")
  must(homeView.Render(w, nil))
}

func contact(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html")
  must(contactView.Render(w, nil))
}

func must(err error) {
  if err != nil {
    panic(err)
  }
}

func main() {
  homeView = views.NewView("bootstrap", "views/home.gohtml")
  contactView = views.NewView("bootstrap", "views/contact.gohtml")
  usersC := controllers.NewUsers()

  r := mux.NewRouter()
  r.HandleFunc("/", home)
  r.HandleFunc("/contact", contact)
  r.HandleFunc("/signup", usersC.New)
  http.ListenAndServe(":3000", r)
}
