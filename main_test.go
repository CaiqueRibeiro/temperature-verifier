package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/CaiqueRibeiro/cloud-run-challenge/handlers"
	"github.com/CaiqueRibeiro/cloud-run-challenge/models"
)

func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()

	// Teardown

	os.Exit(code)
}

func TestHealthCheckHandler(t *testing.T) {
	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HealthCheckHandler)

	// Call the handler with our request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to parse response body: %v", err)
	}
	if response["status"] != "ok" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			response["status"], "ok")
	}
}

func TestWeatherByCEPHandler_InvalidFormat(t *testing.T) {
	// Test with invalid CEP format
	req, err := http.NewRequest("GET", "/weather/cep/123", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.WeatherByCEPHandler)
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnprocessableEntity)
	}

	// Check response body
	var response models.ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to parse response body: %v", err)
	}
	if response.Message != "invalid zipcode" {
		t.Errorf("handler returned unexpected message: got %v want %v",
			response.Message, "invalid zipcode")
	}
}
