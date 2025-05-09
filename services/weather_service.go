package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/CaiqueRibeiro/cloud-run-challenge/models"
)

var (
	weatherAPIBaseURL = "http://api.weatherapi.com/v1"
	weatherAPIKey     = "9a57eaf6223b4d55a8d122153250905"
)

// SetWeatherAPIBaseURL sets the base URL for Weather API (used in tests)
func SetWeatherAPIBaseURL(url string) {
	weatherAPIBaseURL = url
}

// SetWeatherAPIKey sets the API key for Weather API (used in tests)
func SetWeatherAPIKey(key string) {
	weatherAPIKey = key
}

// GetWeatherByCity retrieves weather information for a city
func GetWeatherByCity(city, state string) (*models.WeatherResponse, error) {
	// URL encode the city and state parameters
	city = url.QueryEscape(city)
	state = url.QueryEscape(state)

	// Call WeatherAPI
	url := fmt.Sprintf("%s/current.json?key=%s&q=%s,%s&aqi=no", weatherAPIBaseURL, weatherAPIKey, city, state)
	log.Printf("Calling WeatherAPI: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error making HTTP request to WeatherAPI: %v", err)
		return nil, fmt.Errorf("error calling WeatherAPI: %w", err)
	}
	defer resp.Body.Close()

	// Log response status
	log.Printf("WeatherAPI response status: %s", resp.Status)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading WeatherAPI response body: %v", err)
		return nil, fmt.Errorf("error reading WeatherAPI response: %w", err)
	}

	log.Printf("WeatherAPI response body: %s", string(body))

	// Check for error response
	if resp.StatusCode != http.StatusOK {
		log.Printf("WeatherAPI returned non-OK status: %d", resp.StatusCode)
		return nil, fmt.Errorf("weather API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	var weatherResponse models.WeatherAPIResponse
	if err := json.Unmarshal(body, &weatherResponse); err != nil {
		log.Printf("Error parsing WeatherAPI response: %v", err)
		return nil, fmt.Errorf("error parsing WeatherAPI response: %w", err)
	}

	// Convert temperatures
	tempC := weatherResponse.Current.TempC
	tempF := (tempC * 9 / 5) + 32
	tempK := tempC + 273.15

	log.Printf("Converted temperatures: tempC=%.2f, tempF=%.2f, tempK=%.2f", tempC, tempF, tempK)

	return &models.WeatherResponse{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}, nil
}

// celsiusToFahrenheit converts temperature from Celsius to Fahrenheit
// F = C * 1.8 + 32
func celsiusToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

// celsiusToKelvin converts temperature from Celsius to Kelvin
// K = C + 273
func celsiusToKelvin(celsius float64) float64 {
	return celsius + 273
}
