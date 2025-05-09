# Weather API by CEP

A Golang API that retrieves weather information based on a Brazilian postal code (CEP).

## Features

- Takes a CEP (Brazilian postal code) as input
- Finds the location based on the CEP
- Returns the current temperature in Celsius, Fahrenheit, and Kelvin

## Requirements

- Go 1.21 or higher
- Docker (optional, for containerized deployment)
- WeatherAPI key (sign up at https://www.weatherapi.com/)

## API Endpoints

### Get Weather by CEP

```
GET /weather/cep/{cep}
```

Parameters:
- `cep`: A valid Brazilian postal code (8 digits)

Success Response (200 OK):
```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

Error Responses:
- 422 Unprocessable Entity: Invalid CEP format
  ```json
  {
    "message": "invalid zipcode"
  }
  ```
- 404 Not Found: CEP not found
  ```json
  {
    "message": "can not find zipcode"
  }
  ```

### Health Check

```
GET /health
```

Success Response (200 OK):
```json
{
  "status": "ok"
}
```

## Local Development

### Run with Go

```bash
go run main.go
```

### Run with Docker

```bash
docker build -t weather-api .
docker run -p 8080:8080 weather-api
```

## Testing

Run tests:
```bash
go test ./...
``` 