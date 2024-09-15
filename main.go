package main

import (
	"log"
	"net/http"

	"cabromiley.classes/db"
	"cabromiley.classes/router"
)

func main() {
	database, err := db.InitDB("./users.db")
	if err != nil {
		log.Fatal("Failed to initialize the database:", err)
	}
	defer database.Close()

	r := router.Router(database)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
