package controllers

import (
  "fmt"
  "github.com/jlindauer/usegolang/views"
  "net/http"
  //"github.com/gorilla/schema"
)

// NewUsers returns a User after creating a View
func NewUsers() *Users {
  return &Users{
    NewView: views.NewView("bootstrap", "users/new"),
  }
}

// New is used to render the form where a user can create a new user account
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
  if err := u.NewView.Render(w, nil); err != nil {
    panic(err)
  }
}

// Create is used to process the signup form when a user tries to create a new user account
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
  var form SignupForm
  if err := parseForm(r, &form); err != nil {
    panic(err)
  }

  fmt.Fprintln(w, "Email is", form.Email)
  fmt.Fprintln(w, "Password is", form.Password)
}

type Users struct {
  NewView *views.View
}

type SignupForm struct {
  Email    string `schema:"email"`
  Password string `schema:"password"`
}
