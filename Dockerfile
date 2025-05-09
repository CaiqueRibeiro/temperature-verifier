FROM golang:1.21 AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o weather-api .

FROM scratch
WORKDIR /app
# Copy CA certificates from the build stage
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/weather-api .
ENV HOST 0.0.0.0
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["./weather-api"]