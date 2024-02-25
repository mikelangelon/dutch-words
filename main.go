package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type TimeValue struct {
	Time  string  `json:"time"`
	Value float64 `json:"value"`
}

type Post struct {
	Date    time.Time
	Title   string
	Content string
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
func main() {
	// Setup
	store := NewStore()
	service := NewService(store)
	handler := NewHandler(service)
	mux := http.NewServeMux()

	// Routing
	mux.HandleFunc("GET /web/", func(w http.ResponseWriter, request *http.Request) {
		enableCors(&w)
		Dashboard([]Word{
			{ID: "1", Dutch: "hond", English: "dog"},
		}).Render(request.Context(), w)
	})
	mux.HandleFunc("GET /web/word/{id}", handler.getWord)
	mux.HandleFunc("GET /web/word/dutch/{text}/", handler.getWorByDutch)
	mux.HandleFunc("GET /web/word", handler.getWords)
	mux.HandleFunc("POST /web/word", handler.createWord)
	mux.HandleFunc("DELETE /web/word/{id}/", handler.deleteWord)

	server := &http.Server{Addr: ":8080", Handler: mux}

	fmt.Println("Server is listening on http://localhost:8080")
	err := server.ListenAndServe()
	if err != nil {
		slog.Error("unexpected error on server", "error", err)
	}
}

func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == "OPTIONS" {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}

		next(w, r)
	}
}
