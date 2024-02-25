package main

import (
	"fmt"
	"github.com/a-h/templ"
	"log/slog"
	"net/http"
)

func main() {
	// Setup
	store := NewStore()
	service := NewService(store)
	handler := NewHandler(service)
	mux := http.NewServeMux()

	component := headerTemplate("World")

	// Routing
	mux.Handle("GET localhost/", templ.Handler(component))
	mux.HandleFunc("GET localhost/word/{id}/", handler.getWord)
	mux.HandleFunc("GET localhost/word/dutch/{text}/", handler.getWorByDutch)
	mux.HandleFunc("GET localhost/word/", handler.getWords)
	mux.HandleFunc("POST localhost/word/", handler.createWord)
	mux.HandleFunc("DELETE localhost/word/{id}/", handler.deleteWord)

	server := &http.Server{Addr: ":8080", Handler: mux}

	fmt.Println("Server is listening on http://localhost:8080")
	err := server.ListenAndServe()
	if err != nil {
		slog.Error("unexpected error on server", "error", err)
	}
}
