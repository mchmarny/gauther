package main

import (
	"fmt"
	"log"
	"net/http"
	"context"

	"github.com/mchmarny/gauther/handlers"

	"github.com/mchmarny/gauther/utils"
)



func main() {

	log.Print("Configuring server...")
	ctx := context.Background()

	// config
	port := utils.MustGetEnv("PORT", "8080")
	url := utils.MustGetEnv("EXTERNAL_URL", fmt.Sprintf("http://localhost:%s", port))

	// Google OAuth
	err := handlers.ConfigureOAuthHandler(ctx, url)
	if err != nil {
		log.Fatalf("Error when initializing OAuth handler: %s", err)
	}

	// Mux
	mux := http.NewServeMux()

	// Templates
	mux.Handle("/", http.FileServer(http.Dir("templates/")))

	// Static
	mux.Handle("/static/", http.StripPrefix("/static/",
	  	http.FileServer(http.Dir("static"))))

	// OAuth handlers
	mux.HandleFunc("/auth/login", handlers.OAuthLoginHandler)
	mux.HandleFunc("/auth/callback", handlers.OAuthCallbackHandler)

	// Server configured
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	log.Printf("Server starting on port %s \n", port)
	log.Fatal(server.ListenAndServe())

}
