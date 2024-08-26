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
		log.Fatal("Failed to connect to the database:", err)
	}
	log.Println("Database connection established")

	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT,
		"email" TEXT
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
	log.Println("Table 'users' ensured to exist")
}

func main() {
	initDB()
	defer db.Close()
	log.Println("Database connection closed on exit")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Index handler to list users
func Index(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling Index request")
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal("Failed to query users:", err)
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			log.Fatal("Failed to scan user row:", err)
		}
		users = append(users, user)
	}
	log.Printf("Retrieved %d users", len(users))

	err = tmpl.ExecuteTemplate(w, "Index", users)
	if err != nil {
		log.Println("Failed to execute template for Index:", err)
	}
}

// Show handler to show a single user
func Show(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	log.Printf("Handling Show request for user ID: %s", id)
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		log.Println("Failed to retrieve user:", err)
	}

	err = tmpl.ExecuteTemplate(w, "Show", user)
	if err != nil {
		log.Println("Failed to execute template for Show:", err)
	}
}

// New handler to render the form for creating a new user
func New(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling New request")
	err := tmpl.ExecuteTemplate(w, "New", nil)
	if err != nil {
		log.Println("Failed to execute template for New:", err)
	}
}

// Edit handler to render the form for editing an existing user
func Edit(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	log.Printf("Handling Edit request for user ID: %s", id)
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		log.Println("Failed to retrieve user for editing:", err)
	}

	err = tmpl.ExecuteTemplate(w, "Edit", user)
	if err != nil {
		log.Println("Failed to execute template for Edit:", err)
	}
}

// Insert handler to insert a new user
func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")
		log.Printf("Handling Insert request - Name: %s, Email: %s", name, email)

		if !IsValidEmail(email) {
			log.Println("Invalid email provided:", email)
			http.Error(w, "Invalid Email", http.StatusBadRequest)
			return
		}

		stmt, err := db.Prepare("INSERT INTO users(name, email) VALUES(?, ?)")
		if err != nil {
			log.Fatal("Failed to prepare insert statement:", err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(name, email)
		if err != nil {
			log.Fatal("Failed to execute insert statement:", err)
		}
		log.Println("User inserted successfully")

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	} else {
		log.Println("Invalid method used for Insert")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Update handler to update an existing user
func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id, _ := strconv.Atoi(r.FormValue("id"))
		name := r.FormValue("name")
		email := r.FormValue("email")
		log.Printf("Handling Update request - ID: %d, Name: %s, Email: %s", id, name, email)

		stmt, err := db.Prepare("UPDATE users SET name = ?, email = ? WHERE id = ?")
		if err != nil {
			log.Fatal("Failed to prepare update statement:", err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(name, email, id)
		if err != nil {
			log.Fatal("Failed to execute update statement:", err)
		}
		log.Println("User updated successfully")

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	} else {
		log.Println("Invalid method used for Update")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Delete handler to delete a user
func Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	log.Printf("Handling Delete request for user ID: %s", id)

	stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		log.Fatal("Failed to prepare delete statement:", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal("Failed to execute delete statement:", err)
	}
	log.Println("User deleted successfully")

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
