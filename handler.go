package main

import (
	"encoding/json"
	"net/http"
)

type handler struct {
	Service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		Service: service,
	}
}

func (s *handler) createWord(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	dutch := req.FormValue("dutch")
	english := req.FormValue("english")

	s.Service.InsertWord(Word{Dutch: dutch, English: english})
	words, err := s.Service.FindAllWords()
	if err != nil {
		// TODO Deal with error
	}
	var ws []Word
	for _, v := range words {
		ws = append(ws, *v)
	}
	WordList(ws).Render(req.Context(), w)
}

func (s *handler) deleteWord(w http.ResponseWriter, req *http.Request) {
	err := s.Service.DeleteWord(req.PathValue("id"))
	if err != nil {
		// TODO To improve error codes
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	words, _ := s.Service.FindAllWords()
	var ws []Word
	for _, v := range words {
		ws = append(ws, *v)
	}
	WordList(ws).Render(req.Context(), w)
}

func (s *handler) getWord(w http.ResponseWriter, req *http.Request) {
	word, err := s.Service.FindWordByID(req.PathValue("id"))
	if err != nil {
		// TODO To improve error codes
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderJSON(w, word)
}

func (s *handler) getWorByDutch(w http.ResponseWriter, req *http.Request) {
	word, err := s.Service.FindWordByDutch(req.PathValue("text"))
	if err != nil {
		// TODO To improve error codes
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderJSON(w, word)
}

func (s *handler) getWords(w http.ResponseWriter, req *http.Request) {
	words, err := s.Service.FindAllWords()
	if err != nil {
		// TODO To improve error codes
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderJSON(w, words)
}
func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
