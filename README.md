![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)

# URL Shortener

A minimal URL shortener built with Go and [chi](https://github.com/go-chi/chi), organized with Clean Architecture and DDD concepts. It exposes a small HTTP API to generate short codes for long URLs and redirect visitors when they hit those codes.

## Features

- Shorten any valid http/https URL into a 6-character alphanumeric code
- Redirect from a short code to the original URL
- In-memory storage, safe for concurrent access (no external database required)
- Domain validation: a link can only be created with a valid URL
- Request logging, panic recovery and request ID middleware (via chi)
- CORS enabled for `http://localhost:3000`

## Requirements

- Go 1.22 or later (uses `math/rand/v2`)

## Project Structure

```
.
├── cmd/api/
│   └── main.go                  # Entry point, wires all layers together
├── internal/
│   ├── domain/                  # Business rules, depends on nothing
│   │   ├── link.go              # Link entity, validated by NewLink
│   │   ├── repository.go        # LinkRepository interface
│   │   └── service/
│   │       └── generator.go     # Short code generator
│   ├── exception/
│   │   └── errors.go            # Domain errors (ErrEmptyURL, ErrInvalidURL, ErrNotFound, ErrCodeTaken...)
│   ├── application/             # Use cases
│   │   ├── shorten.go           # Generates a code, creates and saves the link
│   │   └── resolve.go           # Finds the original URL by code
│   ├── infrastructure/
│   │   └── link_repository.go   # In-memory repository (map + RWMutex)
│   └── interfaces/httpapi/      # HTTP layer
│       ├── handler.go           # Routes and handlers (chi)
│       ├── dto.go               # Request/response bodies
│       └── error.go             # Maps domain errors to HTTP status
├── dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum
```

### Dependency rule

```
interfaces/httpapi -> application -> domain <- infrastructure
```

Dependencies always point to the domain. The domain does not know about HTTP, chi or the storage. Replacing the in-memory repository with a database only affects `infrastructure/`.

## Running the project

```bash
go run ./cmd/api
```

The server starts on port `8080`.

## API Reference

### Shorten a URL

Creates a short code for the given URL.

**Request**

```
POST /api/shorten
Content-Type: application/json

{
  "url": "https://www.google.com"
}
```

**Example**

```bash
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://www.google.com"}'
```

**Success response** — `201 Created`

```json
{
  "data": {
    "code": "aZ3kLm"
  }
}
```

### Redirect from a short code

Redirects the client to the original URL associated with a short code.

**Request**

```
GET /{code}
```

**Example**

```bash
curl -i http://localhost:8080/aZ3kLm
```

**Success response** — `302 Found`

Redirects to the original URL (e.g. `https://www.google.com`).

### Errors

All errors use the same body:

```json
{
  "error": "invalid url"
}
```

| Status