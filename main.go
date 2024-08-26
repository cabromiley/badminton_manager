package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

var tmpl = template.Must(template.ParseGlob("templates/*"))

// User represents the user model
type User struct {
	ID    int
	Name  string
	Email string
}

// InitDB initializes the database connection and returns a *sql.DB
func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT,
		"email" TEXT
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}

func main() {
	db, err := InitDB("./users.db")
	if err != nil {
		log.Fatal("Failed to initialize the database:", err)
	}
	defer db.Close()

	// Pass the database connection to your handlers
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", withDB(Index, db))
	http.HandleFunc("/show", withDB(Show, db))
	http.HandleFunc("/new", withDB(New, db))
	http.HandleFunc("/edit", withDB(Edit, db))
	http.HandleFunc("/insert", withDB(Insert, db))
	http.HandleFunc("/update", withDB(Update, db))
	http.HandleFunc("/delete", withDB(Delete, db))

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// withDB is a helper function that injects the db connection into the handler functions
func withDB(handler func(http.ResponseWriter, *http.Request, *sql.DB), db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, db)
	}
}

// Index handler to list users
func Index(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
func Show(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := r.URL.Query().Get("id")
	log.Printf("Handling Show request for user ID: %s", id)
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		log.Println("Failed to retrieve user:", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	err = tmpl.ExecuteTemplate(w, "Show", user)
	if err != nil {
		log.Println("Failed to execute template for Show:", err)
	}
}

// New handler to render the form for creating a new user
func New(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("Handling New request")
	err := tmpl.ExecuteTemplate(w, "New", nil)
	if err != nil {
		log.Println("Failed to execute template for New:", err)
	}
}

// Edit handler to render the form for editing an existing user
func Edit(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := r.URL.Query().Get("id")
	log.Printf("Handling Edit request for user ID: %s", id)
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		log.Println("Failed to retrieve user for editing:", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	err = tmpl.ExecuteTemplate(w, "Edit", user)
	if err != nil {
		log.Println("Failed to execute template for Edit:", err)
	}
}

// Insert handler to insert a new user
func Insert(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
func Update(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
func Delete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
