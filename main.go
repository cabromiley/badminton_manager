package main

import (
	"log"
	"net/http"

	"cabromiley.classes/db"
	"cabromiley.classes/handlers"
	"cabromiley.classes/middleware"
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
	r.HandleFunc("/", middleware.AuthMiddleware(middleware.WithDB(handlers.Index, database)))
	r.HandleFunc("/user/{id}", middleware.WithDB(handlers.ShowUserByID, database)) // Route with dynamic ID
	r.HandleFunc("/new", middleware.WithDB(handlers.New, database))
	r.HandleFunc("/edit/{id}", middleware.WithDB(handlers.Edit, database))
	r.HandleFunc("/insert", middleware.WithDB(handlers.Insert, database))
	r.HandleFunc("/update", middleware.WithDB(handlers.Update, database))
	r.HandleFunc("/delete/{id}", middleware.WithDB(handlers.Delete, database))
	r.HandleFunc("/register", middleware.WithDB(handlers.Register, database))
	r.HandleFunc("/login", middleware.WithDB(handlers.Login, database))
	r.HandleFunc("/logout", handlers.Logout).Methods("GET")

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
