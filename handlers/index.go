package handlers

import (
	"net/http"
	"html/template"
	"log"

	"github.com/mchmarny/gauther/stores"

)


// DefaultHandler handles index page
func DefaultHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Index handler...")

	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	var data map[string]interface{}

	uidCookie, _ := r.Cookie(userIDCookieName)
	if uidCookie != nil && uidCookie.Value != "" {
		log.Printf("User authenticated: %s, getting data...", uidCookie.Value)
		userData, err := stores.GetData(r.Context(), uidCookie.Value)
		if err != nil {
			log.Printf("Error while getting user data: %v", err)
			ErrorHandler(w, r)
			return
		}
		data = userData
	}

	tmpl.Execute(w, data)



}