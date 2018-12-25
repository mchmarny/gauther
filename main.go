package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"context"

	"github.com/mchmarny/gauther/handlers"
)

const (
	portEnvToken = "PORT"
	externalURLToken = "EXTERNAL_URL"
)

func main() {

	log.Print("Configuring server...")

	// context for the entire server instance
	ctx := context.Background()

	// port
	port := os.Getenv(portEnvToken)
	if port == "" {
		port = "8080"
	}

	// port
	url := os.Getenv(externalURLToken)
	if url == "" {
		url = fmt.Sprintf("http://localhost:%s", port)
	}

	// Google OAuth
	err := handlers.ConfigureOAuthHandler(ctx, url)
	if err != nil {
		log.Fatal("Error when initializing OAuth handler")
	}

	// Mux
	mux := http.NewServeMux()

	// Templates
	mux.Handle("/", http.FileServer(http.Dir("templates/")))

	// Static
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

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
