package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jlindauer/usegolang/controllers"
	"github.com/jlindauer/usegolang/models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "usegolang_dev"
)

func main() {
	// Create a DB connection string and then use it to
	// create our model services.
	psqlInfo := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d sslmode=disable",
		dbname, user, password, host, port)

	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	// TODO: Simplify this
	defer services.Close()
	services.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery)

	r := mux.NewRouter()

	// static pages
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.FAQ).Methods("GET")

	// user related routes
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	// gallery routes
	r.Handle("/galleries/new", galleriesC.New).Methods("GET")

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
