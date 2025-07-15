# Build stage
FROM golang:1.24.3 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

# Final image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /OLLAMA-GO ./app
# COPY .env ./
EXPOSE 8080
CMD ["./app"]