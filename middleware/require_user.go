package middleware

import (
  "net/http"
  "strings"

  "github.com/jlindauer/usegolang/models"
  "github.com/jlindauer/usegolang/context"
)

// User middleware will lookup the current user via their
// remember_token cookie using the UserService. If the user
// is found, they will be set on the request context.
// Regardless, the next handler is always called.
type User struct {
  models.UserService
}

func (mw *User) Apply(next http.Handler) http.HandlerFunc {
  return mw.ApplyFn(next.ServeHTTP)
}

func (mw *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path
    // If the user is requesting a static asset or image we will
    // not need to lookup the current user so we skip doing that.
    if strings.HasPrefix(path, "/assets/") || strings.HasPrefix(path, "/images/") {
      next(w, r)
      return
    }

    cookie, err := r.Cookie("remember_token")
    if err != nil {
      next(w, r)
      return
    }
    user, err := mw.UserService.ByRemember(cookie.Value)
    if err != nil {
      next(w, r)
      return
    }
    ctx := r.Context()
    ctx = context.WithUser(ctx, user)
    r = r.WithContext(ctx)
    next(w, r)
  })
}

type RequireUser struct {}

// ApplyFn will return an http.HandlerFunc that will
// check to see if a user i slogged in and then either
// call next(w, r) if they are, or redirect them to the
// login page if they are not
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    user := context.User(r.Context())
    if user == nil {
      http.Redirect(w, r, "/login", http.StatusFound)
      return
    }
    next(w, r)
  })
}

func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
  return mw.ApplyFn(next.ServeHTTP)
}
