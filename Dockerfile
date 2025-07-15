FROM golang:1.24 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/server .

# Final minimal image
FROM gcr.io/distroless/static
COPY --from=builder /app/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]