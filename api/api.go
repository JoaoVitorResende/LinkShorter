package api

import (
	"encoding/json"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

const caracteres = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func NewHandler(db map[string]string) http.Handler {
	r := chi.NewMux()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"},
        AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type"},
        AllowCredentials: false,
        MaxAge:           300,
    }))

	r.Post("/api/shorten", handlePost(db))
	r.Get("/{code}", handleGet(db))
	return r
}

type PostBody struct {
	URL string `json:"url"`
}

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func handlePost(db map[string]string) http.HandlerFunc {
	
	return func(w http.ResponseWriter, r *http.Request) {
		var body PostBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJson(w, Response{Error: "invalid body"}, http.StatusUnprocessableEntity)
			return
		}
		if _, err := url.Parse(body.URL); err != nil {
			sendJson(w, Response{Error: "invalid url passed"}, http.StatusBadRequest)
			return
		}
		code := genCode()
		db[code] = body.URL
		sendJson(w, Response{Data: code}, http.StatusCreated)
	}
}
//curl -X POST http://localhost:8080/api/shorten -d '{"url":"https://www.google.com"}'
func handleGet(db map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := chi.URLParam(r, "code")
		urlDb, ok := db[url]
		if !ok {
			sendJson(w, Response{Error: "Url not found"}, http.StatusNotFound)
			return
		}
		http.Redirect(w, r, urlDb, http.StatusPermanentRedirect)
	}
}

func genCode() string {
	const n = 8
	byts := make([]byte, n)
	for i := range n {
		byts[i] = caracteres[rand.IntN(len(caracteres))]
	}
	return string(byts)
}

func sendJson(w http.ResponseWriter, resp Response, status int) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("failed to marshal json data", "error", err)
		sendJson(
			w,
			Response{Error: "something went wrong"},
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("failed to write json data", "error", err)
		return
	}
}
