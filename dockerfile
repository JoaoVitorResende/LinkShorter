# --- Build stage ---
FROM golang:tip-alpine3.24 AS builder

WORKDIR /app

# Cache de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do codigo
COPY . .

# Build estatico
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server .

# --- Runtime stage ---
FROM alpine:3.24 AS production

WORKDIR /app

# Certificados TLS
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server .

EXPOSE 8080

ENTRYPOINT ["/app/server"]
