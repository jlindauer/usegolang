package main

import (
	"fmt"

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
	psqlInfo := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d sslmode=disable",
		dbname, user, password, host, port)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.DestructiveReset()

	// Create a user
	user := models.User{
		Name:     "Michael Scott",
		Email:    "michael@dundermifflin.com",
		Password: "bestboss",
	}
	if err := us.Create(&user); err != nil {
		panic(err)
	}

	// Verify that the user has a Remember and RememberHash
	fmt.Printf("%+v\n", user)
	if user.Remember == "" {
		panic("Invalid remember token")
	}

	// Now verify that we can lookup a user with that remember token
	user2, err := us.ByRemember(user.Remember)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", *user2)
}
