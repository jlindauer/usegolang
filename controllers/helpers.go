package controllers

import (
  "net/http"
  "github.com/gorilla/schema"
)

// parseForm handles the calling of ParseForm and error handling
// It takes in a Request and an interface and then calls
// Decode to pass the parsed data into the destination interface
func parseForm(r *http.Request, dst interface{}) error {
  if err := r.ParseForm(); err != nil {
    panic(err)
  }

  dec := schema.NewDecoder()
  if err := dec.Decode(dst, r.PostForm); err != nil {
    return err
  }
  return nil
}
