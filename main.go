package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
func main() {
	// Setup
	store := NewStore()
	store.Insert(Word{Dutch: "hond", English: "dog"})
	store.Insert(Word{Dutch: "paard", English: "horse"})
	service := NewService(store)
	handler := NewHandler(service)
	mux := http.NewServeMux()

	// Routing
	mux.HandleFunc("GET /", handler.formAndList)
	mux.HandleFunc("GET /web/", handler.formAndList)
	mux.HandleFunc("GET /web/tab1", handler.tab1)
	mux.HandleFunc("GET /web/tab2", func(w http.ResponseWriter, request *http.Request) {
		Tabs([]Word{}).Render(request.Context(), w)
	})
	mux.HandleFunc("GET /web/word/{id}", handler.getWord)
	mux.HandleFunc("GET /web/word/dutch/{text}/", handler.getWorByDutch)
	mux.HandleFunc("GET /web/word", handler.getWords)
	mux.HandleFunc("POST /web/word", handler.createWord)
	mux.HandleFunc("DELETE /web/word/{id}/", handler.deleteWord)

	server := &http.Server{Addr: "localhost:8080", Handler: mux}

	fmt.Println("Server is listening on http://localhost:8080")
	err := server.ListenAndServe()
	if err != nil {
		slog.Error("unexpected error on server", "error", err)
	}
}
