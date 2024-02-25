package main

import (
	"encoding/json"
	"fmt"
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

func (s *handler) formAndList(w http.ResponseWriter, request *http.Request) {
	enableCors(&w)
	words, _ := s.Service.FindAllWords()
	var ws []Word
	for _, v := range words {
		ws = append(ws, *v)
	}
	Dashboard(ws).Render(request.Context(), w)
}

func (s *handler) tab1(w http.ResponseWriter, request *http.Request) {
	fmt.Println("tab1")
	enableCors(&w)
	words, _ := s.Service.FindAllWords()
	var ws []Word
	for _, v := range words {
		ws = append(ws, *v)
	}
	Tabs(ws).Render(request.Context(), w)
}

func (s *handler) createWord(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	dutch := req.FormValue("dutch")
	english := req.FormValue("english")
	tags := req.Form["tags"]
	wordType := req.FormValue("type")

	s.Service.InsertWord(Word{Dutch: dutch, English: english, Tags: tags, Type: wordType})
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
	search := req.FormValue("word")
	word, err := s.Service.FindWordByDutch(search)
	if err != nil {
		// TODO To improve error codes
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WordList([]Word{*word}).Render(req.Context(), w)
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
