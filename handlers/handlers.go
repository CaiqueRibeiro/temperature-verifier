package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/CaiqueRibeiro/cloud-run-challenge/models"
	"github.com/CaiqueRibeiro/cloud-run-challenge/services"
)

// WeatherByCEPHandler handles the request to get weather information by CEP
func WeatherByCEPHandler(w http.ResponseWriter, r *http.Request) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Extract CEP from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		respondWithError(w, http.StatusBadRequest, "missing CEP parameter")
		return
	}
	cep := pathParts[3]
	if cep == "" {
		respondWithError(w, http.StatusBadRequest, "missing CEP parameter")
		return
	}

	log.Printf("Processing request for CEP: %s", cep)

	// Get location from CEP
	location, err := services.GetLocationByCEP(cep)
	if err != nil {
		log.Printf("Error getting location: %v", err)
		if err.Error() == "invalid zipcode" {
			respondWithError(w, http.StatusUnprocessableEntity, "invalid zipcode")
			return
		}
		if err.Error() == "can not find zipcode" {
			respondWithError(w, http.StatusNotFound, "can not find zipcode")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "error retrieving location")
		return
	}

	log.Printf("Got location: %s, %s", location.Localidade, location.UF)

	// Get weather by city
	weather, err := services.GetWeatherByCity(location.Localidade, location.UF)
	if err != nil {
		log.Printf("Error getting weather: %v", err)
		respondWithError(w, http.StatusInternalServerError, "error retrieving weather information")
		return
	}

	log.Printf("Got weather: tempC=%.2f, tempF=%.2f, tempK=%.2f", weather.TempC, weather.TempF, weather.TempK)

	// Respond with weather information
	respondWithJSON(w, http.StatusOK, weather)
}

// HealthCheckHandler handles the health check endpoint
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	errorResponse := models.ErrorResponse{
		Message: message,
	}
	respondWithJSON(w, code, errorResponse)
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error marshalling JSON response"}`))
		return
	}

	w.WriteHeader(code)
	w.Write(response)
}
