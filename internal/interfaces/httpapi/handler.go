package httpapi

import (
	"encoding/json"
	"net/http"

	"shortapi/internal/application"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Handler(shorten *application.ShortenUseCase, resolve *application.ResolveUseCase) http.Handler {
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

	r.Post("/api/shorten", Post(shorten))
	r.Get("/{code}", Get(resolve))
	return r
}

func Get(resolve *application.ResolveUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		originalURL, err := resolve.Execute(code)
		if err != nil {
			handleError(w, err)
			return
		}
		http.Redirect(w, r, originalURL, http.StatusFound)
	}
}

func Post(shorten *application.ShortenUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body PostBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			handleError(w, errInvalidBody)
			return
		}
		code, err := shorten.Execute(body.URL)
		if err != nil {
			handleError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, Response{Data: map[string]string{"code": code}})
	}
}