.PHONY: build run test docker-build docker-run help

help:
	@echo "Available commands:"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application locally"
	@echo "  make test          - Run tests"
	@echo "  make docker-build  - Build the Docker image"
	@echo "  make docker-run    - Run the application in Docker"
	@echo "  make help          - Show this help message"

build:
	go build -o weather-api main.go

run:
	go run main.go

test:
	go test ./...

docker-build:
	docker build -t weather-api .

docker-run:
	docker run -p 8080:8080 -e WEATHER_API_KEY=${WEATHER_API_KEY} weather-api 