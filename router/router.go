package router

import (
	"database/sql"
	"net/http"

	"cabromiley.classes/handlers"
	"cabromiley.classes/middleware"
	"github.com/gorilla/mux"
)

func Router(db *sql.DB) *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/", middleware.AuthMiddleware(middleware.WithDB(handlers.Index, db)))
	r.HandleFunc("/user/{id}", middleware.WithDB(handlers.ShowUserByID, db)) // Route with dynamic ID
	r.HandleFunc("/new", middleware.WithDB(handlers.New, db))
	r.HandleFunc("/edit/{id}", middleware.WithDB(handlers.Edit, db))
	r.HandleFunc("/insert", middleware.WithDB(handlers.Insert, db))
	r.HandleFunc("/update", middleware.WithDB(handlers.Update, db))
	r.HandleFunc("/delete/{id}", middleware.WithDB(handlers.Delete, db))
	r.HandleFunc("/register", middleware.WithDB(handlers.Register, db))
	r.HandleFunc("/login", middleware.WithDB(handlers.Login, db))
	r.HandleFunc("/logout", handlers.Logout).Methods("GET")

	return r
}
