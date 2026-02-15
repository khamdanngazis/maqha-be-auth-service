# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/auth_service ./cmd/main.go

# Runtime stage
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /bin/auth_service /app/auth_service
COPY cmd/config /app/cmd/config

RUN mkdir -p /app/logs

EXPOSE 8011 50053

ENTRYPOINT ["/app/auth_service"]
CMD ["-config", "cmd/config/config.yaml", "-log.file", "logs"]
