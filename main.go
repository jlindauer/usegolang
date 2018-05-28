package main

import (
  "net/http"
  "github.com/gorilla/mux"
  "github.com/jlindauer/usegolang/controllers"
)

func main() {
  staticC := controllers.NewStatic()
  usersC := controllers.NewUsers()

  r := mux.NewRouter()
  r.Handle("/", staticC.Home).Methods("GET")
  r.Handle("/contact", staticC.Contact).Methods("GET")
  r.HandleFunc("/signup", usersC.New).Methods("GET")
  r.HandleFunc("/signup", usersC.Create).Methods("POST")
  http.ListenAndServe(":3000", r)
}
