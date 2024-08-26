package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var tmpl = template.Must(template.ParseGlob("templates/*"))

// User represents the user model
type User struct {
	ID    int
	Name  string
	Email string
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT,
		"email" TEXT
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initDB()
	defer db.Close()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

// Index handler to list users
func Index(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	tmpl.ExecuteTemplate(w, "Index", users)
}

// Show handler to show a single user
func Show(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		log.Fatal(err)
	}

	tmpl.ExecuteTemplate(w, "Show", user)
}

// New handler to render the form for creating a new user
func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

// Edit handler to render the form for editing an existing user
func Edit(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		log.Fatal(err)
	}

	tmpl.ExecuteTemplate(w, "Edit", user)
}

// Insert handler to insert a new user
func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")

		if !IsValidEmail(email) {
			http.Error(w, "Invalid Email passed", http.StatusBadRequest)
			return
		}

		stmt, err := db.Prepare("INSERT INTO users(name, email) VALUES(?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(name, email)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Update handler to update an existing user
func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id, _ := strconv.Atoi(r.FormValue("id"))
		name := r.FormValue("name")
		email := r.FormValue("email")

		stmt, err := db.Prepare("UPDATE users SET name = ?, email = ? WHERE id = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(name, email, id)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/", 301)
	}
}

// Delete handler to delete a user
func Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/", 301)
}
