package handlers

import (
	"net/http"
	"log"

	"github.com/mchmarny/gauther/stores"
	"github.com/mchmarny/gauther/utils"

)


// DefaultHandler handles index page
func DefaultHandler(w http.ResponseWriter, r *http.Request) {

	var data map[string]interface{}
	data["version"] = utils.MustGetEnv("RELEASE", "NOT SET")

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

	if err := templates.ExecuteTemplate(w, "home", data); err != nil {
		log.Printf("Error in home template: %s", err)
	}



}