package handlers

import (
	"log"
	"net/http"
)

// ErrorHandler handles all app errors
func ErrorHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Error handler not implemented yet...")
	http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)

}
