package services

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLocationByCEP(t *testing.T) {
	tests := []struct {
		name          string
		cep           string
		mockResponse  string
		mockStatus    int
		expectedError error
	}{
		{
			name:          "invalid CEP format",
			cep:           "123",
			expectedError: errors.New("invalid zipcode"),
		},
		{
			name:         "valid CEP",
			cep:          "01001000",
			mockResponse: `{"cep": "01001-000", "logradouro": "Praça da Sé", "complemento": "lado ímpar", "bairro": "Sé", "localidade": "São Paulo", "uf": "SP", "ibge": "3550308", "gia": "1004", "ddd": "11", "siafi": "7107"}`,
			mockStatus:   http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// If we have a mock response, set up a test server
			if tt.mockResponse != "" {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.mockStatus)
					w.Write([]byte(tt.mockResponse))
				}))
				defer server.Close()
			}

			location, err := GetLocationByCEP(tt.cep)

			// Check error cases
			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.expectedError)
					return
				}
				if err.Error() != tt.expectedError.Error() {
					t.Errorf("Expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			// For valid cases
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if location == nil {
				t.Errorf("Expected location data, got nil")
				return
			}
		})
	}
}
