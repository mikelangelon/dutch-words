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

func (s *handler) createWord(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	dutch := req.FormValue("dutch")
	english := req.FormValue("english")

	type ResponseID struct {
		ID string `json:"id"`
	}
	fmt.Println(fmt.Sprintf("%s %s", dutch, english))
	renderJSON(w, ResponseID{ID: "1"})
}

func (s *handler) deleteWord(w http.ResponseWriter, req *http.Request) {
	err := s.Service.DeleteWord(req.PathValue("id"))
	if err != nil {
		// TODO To improve error codes
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
