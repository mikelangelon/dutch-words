package server

import (
	"github.com/mikelangelon/dutch-words/components"
	"github.com/mikelangelon/dutch-words/services"
	"net/http"
)

func New(service services.Service) *http.Server {
	handler := newHandler(service)
	mux := http.NewServeMux()

	// Routing
	mux.HandleFunc("GET /", handler.formAndList)
	mux.HandleFunc("GET /web/", handler.formAndList)
	mux.HandleFunc("GET /web/tab1", handler.tab1)
	mux.HandleFunc("GET /web/tab2", func(w http.ResponseWriter, request *http.Request) {
		navBar := components.NavBar(nav("Search"))
		components.Dashboard(navBar, components.WordSearch()).Render(request.Context(), w)
	})
	mux.HandleFunc("GET /web/tab3", func(w http.ResponseWriter, request *http.Request) {
		navBar := components.NavBar(nav("Tags"))
		components.Dashboard(navBar, components.Tags([]string{"verb", "adjective", "animal"})).Render(request.Context(), w)
	})
	mux.HandleFunc("GET /web/tags/{tag}", handler.tags)
	mux.HandleFunc("GET /web/word/{id}", handler.getWord)
	mux.HandleFunc("GET /web/word/dutch/{text}/", handler.getWorByDutch)
	mux.HandleFunc("GET /web/word", handler.getWords)
	mux.HandleFunc("POST /web/word", handler.createWord)
	mux.HandleFunc("DELETE /web/word/{id}", handler.deleteWord)

	return &http.Server{Addr: "localhost:8080", Handler: mux}
}
