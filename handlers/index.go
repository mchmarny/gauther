package handlers

import (
	"net/http"
	"log"

	"github.com/mchmarny/gauther/stores"

)


// DefaultHandler handles index page
func DefaultHandler(w http.ResponseWriter, r *http.Request) {

	var data map[string]interface{}

	uid := getCurrentUserID(r)
	if uid != "" {
		log.Printf("User authenticated: %s, getting data...", uid)
		userData, err := stores.GetData(r.Context(), uid)
		if err != nil {
			log.Printf("Error while getting user data: %v", err)
			ErrorHandler(w, r, err, http.StatusInternalServerError)
			return
		}
		data = userData
	}

	templates.ExecuteTemplate(w, "home", data)



}