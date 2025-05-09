package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/CaiqueRibeiro/cloud-run-challenge/models"
)

var (
	viaCEPBaseURL = "https://viacep.com.br"
)

// SetViaCEPBaseURL sets the base URL for ViaCEP API (used in tests)
func SetViaCEPBaseURL(url string) {
	viaCEPBaseURL = url
}

// GetLocationByCEP retrieves location information from a CEP (Brazilian postal code)
func GetLocationByCEP(cep string) (*models.ViaCEPResponse, error) {
	// Remove any non-digit characters
	cep = regexp.MustCompile(`[^0-9]`).ReplaceAllString(cep, "")

	// Validate CEP format (8 digits)
	if len(cep) != 8 {
		return nil, errors.New("invalid zipcode")
	}

	// Call ViaCEP API
	url := fmt.Sprintf("%s/ws/%s/json/", viaCEPBaseURL, cep)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error calling ViaCEP API: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading ViaCEP response: %w", err)
	}

	// Parse JSON response
	var cepResponse models.ViaCEPResponse
	if err := json.Unmarshal(body, &cepResponse); err != nil {
		return nil, fmt.Errorf("error parsing ViaCEP response: %w", err)
	}

	// Check if CEP was found
	if cepResponse.Erro == "true" || cepResponse.CEP == "" {
		return nil, errors.New("can not find zipcode")
	}

	return &cepResponse, nil
}
