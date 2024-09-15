package middleware

import (
	"database/sql"
	"net/http"
)

func WithDB(handler func(http.ResponseWriter, *http.Request, *sql.DB), db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, db)
	}
}
