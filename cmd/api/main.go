package main

import (
	"log"
	"net/http"

	"shortapi/internal/application"
	codegen "shortapi/internal/domain/service"
	memory "shortapi/internal/infrastructure"
	"shortapi/internal/interfaces/httpapi"
)

func main() {
	repo := memory.NewLinkRepository()
	gen := codegen.New(6)

	shorten := application.NewShortenUseCase(repo, gen)
	resolve := application.NewResolveUseCase(repo)

	handler := httpapi.Handler(shorten, resolve)

	log.Println("servidor em :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}