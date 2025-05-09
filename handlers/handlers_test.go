package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CaiqueRibeiro/cloud-run-challenge/models"
)

func TestHealthCheckHandler(t *testing.T) {
	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)

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

func TestWeatherByCEPHandler_Integration(t *testing.T) {
	// Skip this test in short mode as it calls real external APIs
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create test cases
	tests := []struct {
		name           string
		cep            string
		expectedStatus int
		expectedBody   models.ErrorResponse
	}{
		{
			name:           "missing CEP",
			cep:            "",
			expectedStatus: http.StatusBadRequest,
			expectedBody: models.ErrorResponse{
				Message: "missing CEP parameter",
			},
		},
		{
			name:           "invalid CEP format",
			cep:            "123",
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody: models.ErrorResponse{
				Message: "invalid zipcode",
			},
		},
		{
			name:           "CEP not found",
			cep:            "99999999",
			expectedStatus: http.StatusNotFound,
			expectedBody: models.ErrorResponse{
				Message: "can not find zipcode",
			},
		},
		{
			name:           "valid CEP",
			cep:            "01001000", // São Paulo, SP
			expectedStatus: http.StatusOK,
			expectedBody:   models.ErrorResponse{},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest("GET", "/weather/cep/"+tt.cep, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler := http.HandlerFunc(WeatherByCEPHandler)
			handler.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			// Check response body
			if tt.expectedStatus != http.StatusOK {
				var response models.ErrorResponse
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatal(err)
				}
				if response.Message != tt.expectedBody.Message {
					t.Errorf("handler returned unexpected message: got %v want %v", response.Message, tt.expectedBody.Message)
				}
			} else {
				// For successful response, check temperature values
				var response models.WeatherResponse
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatal(err)
				}

				// Just check that we have valid temperature values
				if response.TempC == 0 {
					t.Error("expected temperature in Celsius to be non-zero")
				}
				if response.TempF == 0 {
					t.Error("expected temperature in Fahrenheit to be non-zero")
				}
				if response.TempK == 0 {
					t.Error("expected temperature in Kelvin to be non-zero")
				}
			}
		})
	}
}
