package main

import (
	"log"
	"net/http"

	"cabromiley.classes/db"
	"cabromiley.classes/handlers"
	"github.com/gorilla/mux"
)

func main() {
	database, err := db.InitDB("./users.db")
	if err != nil {
		log.Fatal("Failed to initialize the database:", err)
	}
	defer database.Close()

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/", handlers.AuthMiddleware(handlers.WithDB(handlers.Index, database)))
	r.HandleFunc("/user/{id}", handlers.WithDB(handlers.ShowUserByID, database)) // Route with dynamic ID
	r.HandleFunc("/new", handlers.WithDB(handlers.New, database))
	r.HandleFunc("/edit/{id}", handlers.WithDB(handlers.Edit, database))
	r.HandleFunc("/insert", handlers.WithDB(handlers.Insert, database))
	r.HandleFunc("/update", handlers.WithDB(handlers.Update, database))
	r.HandleFunc("/delete/{id}", handlers.WithDB(handlers.Delete, database))
	r.HandleFunc("/register", handlers.WithDB(handlers.Register, database))
	r.HandleFunc("/login", handlers.WithDB(handlers.Login, database))

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
