package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mchmarny/gauther/handlers"
)

const (
	defaultPort = "8080"
	portEnvName = "PORT"
)

func main() {

	log.Print("Configuring server...")

	// port
	port := os.Getenv(portEnvName)
	if port == "" {
		port = defaultPort
	}

	mux := http.NewServeMux()

	// Templates
	mux.Handle("/", http.FileServer(http.Dir("templates/")))

	// Static
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Google Oauth
	handlers.ConfigureOAuthProvider("http://localhost:8080")
	mux.HandleFunc("/auth/login", handlers.OAuthLoginHandler)
	mux.HandleFunc("/auth/callback", handlers.OAuthCallbackHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	log.Printf("Server starting on port %s \n", port)
	log.Fatal(server.ListenAndServe())

}
