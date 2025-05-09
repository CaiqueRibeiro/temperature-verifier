package services

import (
	"testing"
)

func TestTemperatureConversions(t *testing.T) {
	tests := []struct {
		name           string
		celsius        float64
		wantFahrenheit float64
		wantKelvin     float64
	}{
		{
			name:           "freezing point",
			celsius:        0,
			wantFahrenheit: 32,
			wantKelvin:     273,
		},
		{
			name:           "boiling point",
			celsius:        100,
			wantFahrenheit: 212,
			wantKelvin:     373,
		},
		{
			name:           "room temperature",
			celsius:        20,
			wantFahrenheit: 68,
			wantKelvin:     293,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFahrenheit := celsiusToFahrenheit(tt.celsius)
			if gotFahrenheit != tt.wantFahrenheit {
				t.Errorf("celsiusToFahrenheit() = %v, want %v", gotFahrenheit, tt.wantFahrenheit)
			}

			gotKelvin := celsiusToKelvin(tt.celsius)
			if gotKelvin != tt.wantKelvin {
				t.Errorf("celsiusToKelvin() = %v, want %v", gotKelvin, tt.wantKelvin)
			}
		})
	}
}
