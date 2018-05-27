package controllers

import (
  "github.com/jlindauer/usegolang/views"
  "net/http"
)

func NewUsers() *Users {
  return &Users{
    NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
  }
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
  if err := u.NewView.Render(w, nil); err != nil {
    panic(err)
  }
}

type Users struct {
  NewView *views.View
}
