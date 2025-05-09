package main

import (
	"log"
	"net/http"
	"os"

	"github.com/CaiqueRibeiro/cloud-run-challenge/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Set up routes
	http.HandleFunc("/weather/cep/", handlers.WeatherByCEPHandler)
	http.HandleFunc("/health", handlers.HealthCheckHandler)

	// Start server
	addr := ":" + port
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
