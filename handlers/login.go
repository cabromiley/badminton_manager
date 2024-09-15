package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"cabromiley.classes/middleware"
	"cabromiley.classes/utils"
	"golang.org/x/crypto/bcrypt"
)

// Login handler - serves the login form and handles authentication
func Login(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == "GET" {
		// Serve the login form
		err := tmpl.ExecuteTemplate(w, "layout.html", map[string]interface{}{
			"Page": "Login",
		})
		if err != nil {
			log.Println("Failed to render template for login form:", err)
		}
	} else if r.Method == "POST" {
		// Handle form submission (login attempt)
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Basic validation
		if email == "" || password == "" {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		if !utils.IsValidEmail(email) {
			http.Error(w, "Invalid email format", http.StatusBadRequest)
			return
		}

		// Query the user by email
		var hashedPassword string
		var name string
		err := db.QueryRow("SELECT name, password FROM users WHERE email = ?", email).Scan(&name, &hashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			} else {
				log.Println("Error querying user:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// Compare the hashed password with the provided password
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Create a session and store user information
		session, _ := middleware.Store.Get(r, "session")
		session.Values["authenticated"] = true
		session.Values["user"] = name
		err = session.Save(r, w)
		if err != nil {
			log.Println("Failed to save session:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Redirect to the home page (or a protected page)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Logout handler - clears the session and logs the user out
func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.Store.Get(r, "session")
	// Set the authenticated value to false and save the session
	session.Values["authenticated"] = false
	session.Save(r, w)

	// Redirect to the login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Register handler - serves the registration form and handles form submission
func Register(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == "GET" {
		// Serve the registration form
		if r.Header.Get("HX-Request") != "" {
			err := tmpl.ExecuteTemplate(w, "Register.html", map[string]interface{}{
				"Page": "Register",
			})

			if err != nil {
				log.Println("Failed to render template for registration form:", err)
			}
		} else {
			err := tmpl.ExecuteTemplate(w, "layout.html", map[string]interface{}{
				"Page": "Register",
			})

			if err != nil {
				log.Println("Failed to render template for registration form:", err)
			}
		}

	} else if r.Method == "POST" {
		// Handle form submission
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Basic validation
		if name == "" || email == "" || password == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		if !utils.IsValidEmail(email) {
			http.Error(w, "Invalid email format", http.StatusBadRequest)
			return
		}

		// Hash the password using bcrypt
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Failed to hash password:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Insert the user into the database
		stmt, err := db.Prepare("INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)")
		if err != nil {
			log.Println("Failed to prepare statement for inserting user:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(name, email, hashedPassword, "unverified")
		if err != nil {
			log.Println("Failed to insert user into the database:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Redirect to the login page or show success message
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
