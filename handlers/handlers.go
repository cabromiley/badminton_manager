package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"cabromiley.classes/models"
	"cabromiley.classes/utils"
	"github.com/gorilla/mux"
)

var tmpl = template.Must(template.ParseGlob("templates/*"))

// WithDB is a helper function that injects the db connection into the handler functions
func WithDB(handler func(http.ResponseWriter, *http.Request, *sql.DB), db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, db)
	}
}

func AuthMiddleware(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			// If the user is not authenticated, redirect to the login page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		handler(w, r)
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

	users := []models.User{}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
		if err != nil {
			log.Fatal("Failed to scan user row:", err)
		}
		users = append(users, user)
	}
	log.Printf("Retrieved %d users", len(users))

	if r.Header.Get("HX-Request") != "" {
		// Render only the edit content
		err = tmpl.ExecuteTemplate(w, "Index.html", map[string]interface{}{
			"Users": users,
			"Page":  "index",
		})
	} else {
		err = tmpl.ExecuteTemplate(w, "layout.html", map[string]interface{}{
			"Users": users,
			"Page":  "index",
		})
	}
	if err != nil {
		log.Println("Failed to execute template for Index:", err)
	}
}

// ShowUserByID handler to show a single user by ID extracted from the URL path
func ShowUserByID(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := mux.Vars(r)["id"]
	log.Printf("Handling Show request for user ID: %s", id)
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		log.Println("Failed to retrieve user:", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if r.Header.Get("HX-Request") != "" {
		// Render only the edit content
		err = tmpl.ExecuteTemplate(w, "Show.html", map[string]interface{}{
			"User": user,
			"Page": "show",
		})
	} else {
		err = tmpl.ExecuteTemplate(w, "layout.html", map[string]interface{}{
			"User": user,
			"Page": "show",
		})
	}

	if err != nil {
		log.Println("Failed to execute template for Show:", err)
	}
}

func New(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("Handling New request")
	if r.Header.Get("HX-Request") != "" {
		// Render only the edit content
		err := tmpl.ExecuteTemplate(w, "New.html", nil)

		if err != nil {
			log.Println("Failed to execute template for New:", err)
		}
	} else {
		err := tmpl.ExecuteTemplate(w, "layout.html", map[string]interface{}{
			"Page": "show",
		})

		if err != nil {
			log.Println("Failed to execute template for New:", err)
		}
	}

}

func Edit(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := mux.Vars(r)["id"]
	log.Printf("Handling Edit request for user ID: %s", id)
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		log.Println("Failed to retrieve user for editing:", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if the request comes from HTMX
	if r.Header.Get("HX-Request") != "" {
		// Render only the edit content
		err = tmpl.ExecuteTemplate(w, "Edit.html", user)
	} else {
		// Render the full layout with the edit content
		err = tmpl.ExecuteTemplate(w, "layout.html", map[string]interface{}{
			"User": user,
			"Page": "edit",
		})
	}
	if err != nil {
		log.Println("Failed to execute template for Edit:", err)
	}
}

func Insert(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")
		log.Printf("Handling Insert request - Name: %s, Email: %s", name, email)

		if !utils.IsValidEmail(email) {
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

		// Redirect to the index page via HTMX
		Index(w, r, db)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

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

		// Redirect to the index page via HTMX
		Index(w, r, db)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Delete handler to delete a user
func Delete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := mux.Vars(r)["id"]
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

	// Redirect to the index page via HTMX
	Index(w, r, db)
}
