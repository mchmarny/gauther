package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mchmarny/gauther/handlers"
)

const (
	defaultPort      = "8080"
	portVariableName = "PORT"
)

func main() {

	port := os.Getenv(portVariableName)
	if port == "" {
		port = defaultPort
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handlers.NewOAuthHandler(),
	}

	log.Printf("Starting server at %q", server.Addr)
	log.Fatal(server.ListenAndServe())

}
