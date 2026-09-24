![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)

# URL Shortener

A minimal URL shortener built with Go and [chi](https://github.com/go-chi/chi). It exposes a small HTTP API to generate short codes for long URLs and redirect visitors when they hit those codes.

## Features

- Shorten any valid URL into an 8-character alphanumeric code
- Redirect from a short code to the original URL
- In-memory storage (no external database required)
- Structured logging via `log/slog`
- Request logging, panic recovery and request ID middleware (via chi)

## Requirements

- Go 1.22 or later (uses `for i := range n` integer range and `math/rand/v2`)

## Project Structure

```
.
├── main.go          # Entry point, HTTP server setup
└── api/
    └── handler.go    # Routes, handlers and code generation logic
```

## Running the project

```bash
go run main.go
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
  -d '{"url":"https://www.google.com"}'
```

**Success response** — `201 Created`

```json
{
  "data": "aZ3kLm9Q"
}
```

**Error response** — `422 Unprocessable Entity` (invalid JSON body)

```json
{
  "error": "invalid body"
}
```

**Error response** — `400 Bad Request` (invalid URL)

```json
{
  "error": "invalid url passed"
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
curl -i http://localhost:8080/aZ3kLm9Q
```

**Success response** — `308 Permanent Redirect`

Redirects to the original URL (e.g. `https://www.google.com`).

**Error response** — `404 Not Found`

```json
{
  "error": "Url not found"
}
```

## How code generation works

Each short code is 8 characters long, randomly picked from the set `A-Z`, `a-z`, `0-9` (62 possible characters), giving 62^8 (~218 trillion) possible combinations. See `genCode()` in `api/handler.go`.

## Data storage

URLs are stored in an in-memory `map[string]string` (`code -> original URL`). This means:

- All data is lost when the server restarts
- The map is not safe for concurrent access without external synchronization

For production use, consider replacing the in-memory map with a persistent store (e.g. PostgreSQL, Redis) and adding a mutex or using a concurrent-safe map if the in-memory version is kept for testing.

## Known limitations

- No validation preventing duplicate short codes (collision is possible, though unlikely, with 62^8 combinations)
- No uniqueness check for the same URL being shortened multiple times (each request creates a new code)
- The handler for invalid URLs does not stop execution after sending the error response, so processing may continue unexpectedly
- No authentication or rate limiting on the `/api/shorten` endpoint

## License

Not specified.
