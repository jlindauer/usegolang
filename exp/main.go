package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
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
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	var id int
	var name, email string
	rows, err := db.Query(`
    SELECT id, name, email
    FROM users
    WHERE email=$1
		OR ID > $2`,
		"jon@calhoun.io", 3)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		rows.Scan(&id, &name, &email)
		fmt.Println("ID:", id, "Name:", name, "Email:", email)
	}

	db.Close()
}
