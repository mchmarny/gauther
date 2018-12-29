package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// ErrorHandler handles all app errors
func ErrorHandler(w http.ResponseWriter, r *http.Request, err error, code int) {

	log.Printf("Error: %v", err)
	errMsg := fmt.Sprintf("%+v", err)

	tmpl := template.Must(template.ParseFiles("templates/error.html"))

	w.WriteHeader(code)
	tmpl.Execute(w, map[string]interface{}{
		"error":       errMsg,
		"status_code": code,
		"status":      http.StatusText(code),
	})

}
