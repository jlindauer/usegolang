package controllers

import "github.com/jlindauer/usegolang/views"

// NewStatic creates any static pages that do not require a full controller
func NewStatic() *Static {
  return &Static{
    Home: views.NewView("bootstrap", "static/home"),
    Contact: views.NewView("bootstrap", "static/contact"),
  }
}

type Static struct {
  Home *views.View
  Contact *views.View
}
