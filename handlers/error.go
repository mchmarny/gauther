package handlers

import (
	"log"
	"net/http"
)

// ErrorHandler handles index page
func ErrorHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Error handler...")
	http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)

}
