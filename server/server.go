package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/mikelangelon/dutch-words/core"

	"github.com/mikelangelon/dutch-words/components"
	"github.com/mikelangelon/dutch-words/services"
)

func New(service services.Service, ss services.SentencesService) *http.Server {
	handler := newHandler(service, ss)
	mux := http.NewServeMux()

	// Routing
	mux.HandleFunc("GET /", handler.formAndList)
	mux.HandleFunc("GET /web/", handler.formAndList)
	mux.HandleFunc("GET /web/tab1", handler.tab1)
	mux.HandleFunc("GET /web/sentence-tab", func(w http.ResponseWriter, request *http.Request) {
		enableCors(&w)
		navBar := components.NavBar(nav("Sentences"))
		sentences, err := handler.SentencesService.FindAll()
		if err != nil {
			slog.Error("problem retrieving all sentences", "error", err)
			return
		}
		list := formAndList(sentences)
		err = components.Dashboard(navBar, list).Render(request.Context(), w)
		if err != nil {
			slog.Error("problem rendering dashboard", "error", err)
		}
	})
	mux.HandleFunc("GET /web/tab2", func(w http.ResponseWriter, request *http.Request) {
		navBar := components.NavBar(nav("Search"))
		err := components.Dashboard(navBar, components.WordSearch()).Render(request.Context(), w)
		if err != nil {
			slog.Error("problem rendering dashboard", "error", err)
		}
	})
	mux.HandleFunc("GET /web/tab3", handler.renderTagsScreen)
	mux.HandleFunc("POST /web/tags", handler.newTags)
	mux.HandleFunc("GET /web/tags/{tag}", handler.tags)
	mux.HandleFunc("GET /web/word/{id}", handler.getWord)
	mux.HandleFunc("GET /web/word/{id}/edit", handler.getWordEdit)
	mux.HandleFunc("GET /web/word/dutch/{text}/", handler.getWordByDutch)
	mux.HandleFunc("GET /web/word", handler.getWords)
	mux.HandleFunc("POST /web/word", handler.createWord)
	mux.HandleFunc("PUT /web/word/{id}", handler.putWord)
	mux.HandleFunc("DELETE /web/word/{id}", handler.deleteWord)

	mux.HandleFunc("GET /web/sentences/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
		sentence, err := handler.SentencesService.FindById(r.PathValue("id"))
		if err != nil {
			// TODO To improve error codes
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = components.SentenceEdit(core.NewSentenceData(*sentence)).Render(r.Context(), w)
		if err != nil {
			slog.Error("problem rendering word card edit")
		}
	})
	mux.HandleFunc("GET /web/sentences/{id}", func(w http.ResponseWriter, r *http.Request) {
		sentence, err := handler.SentencesService.FindById(r.PathValue("id"))
		if err != nil {
			// TODO To improve error codes
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = components.SentencesCard(*sentence).Render(r.Context(), w)
		if err != nil {
			slog.Error("problem rendering word card edit")
		}
	})
	mux.HandleFunc("POST /web/sentences", func(writer http.ResponseWriter, request *http.Request) {
		sentence := core.Sentence{
			ID:      fmt.Sprintf("%d", time.Now().UnixNano()),
			Dutch:   request.FormValue("dutch"),
			English: request.FormValue("english"),
		}
		if err := handler.SentencesService.Insert(&sentence); err != nil {
			slog.Error("problem inserting sentence", "error", err)
		}
		err := components.SentenceEdit(core.SentenceData{Sentence: core.Sentence{}}).Render(context.TODO(), writer)
		if err != nil {
			return
		}
		err = components.SentenceExtra(sentence).Render(request.Context(), writer)
		if err != nil {
			slog.Error("problem rendering extra sentence", "error", err)
		}
	})
	mux.HandleFunc("PUT /web/sentences/{id}", func(writer http.ResponseWriter, request *http.Request) {
		sentence := &core.Sentence{
			ID:      request.PathValue("id"),
			Dutch:   request.FormValue("dutch"),
			English: request.FormValue("english"),
		}
		if err := handler.SentencesService.Update(sentence); err != nil {
			slog.Error("problem update sentence", "error", err)
		}
		err := components.SentencesCard(*sentence).Render(context.TODO(), writer)
		if err != nil {
			return
		}
	})
	mux.HandleFunc("DELETE /web/sentences/{id}", func(writer http.ResponseWriter, request *http.Request) {
		if err := handler.SentencesService.Delete(request.PathValue("id")); err != nil {
			slog.Error("problem deleting sentence", "error", err)
		}
	})
	mux.HandleFunc("GET /web/sentences", func(writer http.ResponseWriter, request *http.Request) {
		sentences, err := handler.SentencesService.FindAll()
		if err != nil {
			slog.Error("problem getting sentences", "error", err)
			return
		}
		err = components.SentencesList(sentences).Render(context.TODO(), writer)
		if err != nil {
			slog.Error("problem rendering sentences", "error", err)
		}
	})

	return &http.Server{Addr: "localhost:8080", Handler: mux}
}
