package server

import (
	"encoding/json"
	"github.com/mikelangelon/dutch-words/components"
	"github.com/mikelangelon/dutch-words/core"
	"github.com/mikelangelon/dutch-words/services"
	"net/http"
)

type handler struct {
	Service services.Service
}

func newHandler(service services.Service) *handler {
	return &handler{
		Service: service,
	}
}

func (s *handler) tags(w http.ResponseWriter, request *http.Request) {
	tag := request.PathValue("tag")
	search := core.Search{
		Tag: &tag,
	}
	words, err := s.Service.FindWordsBy(search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var ws []core.Word
	for _, v := range words {
		ws = append(ws, *v)
	}
	components.WordList(ws).Render(request.Context(), w)
}
func (s *handler) formAndList(w http.ResponseWriter, request *http.Request) {
	enableCors(&w)
	navBar := components.NavBar(nav("Home"))
	words, _ := s.Service.FindAllWords()
	var ws []core.Word
	for _, v := range words {
		ws = append(ws, *v)
	}
	tab1 := components.Tabs(core.NewFormData(), ws)
	components.Dashboard(navBar, tab1).Render(request.Context(), w)
}

func nav(current string) core.NavigationItems {
	items := core.NavigationItems{
		{
			Label:  "Home",
			Link:   "/web/tab1",
			Active: false,
		},
		{
			Label:  "Search",
			Link:   "/web/tab2",
			Active: false,
		},
		{
			Label:  "Tags",
			Link:   "/web/tab3",
			Active: false,
		},
	}
	for _, v := range items {
		if v.Label == current {
			v.Active = true
		}
	}
	return items
}
func (s *handler) tab1(w http.ResponseWriter, request *http.Request) {
	enableCors(&w)
	words, _ := s.Service.FindAllWords()
	var ws []core.Word
	for _, v := range words {
		ws = append(ws, *v)
	}
	navBar := components.NavBar(nav("Home"))
	components.Dashboard(navBar, components.Tabs(core.NewFormData(), ws)).Render(request.Context(), w)
}

func (s *handler) createWord(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	dutch := req.FormValue("dutch")
	english := req.FormValue("english")
	tags := req.Form["tags"]
	wordType := req.FormValue("type")

	if dutch == "hond" {
		data := core.NewFormData()
		data.Errors["word"] = "something crazy"
		components.WordForm(data).Render(req.Context(), w)
		http.Error(w, "duplicated", http.StatusUnprocessableEntity)
		return
	}
	word := &core.Word{Dutch: dutch, English: english, Tags: tags, Type: wordType}
	err := s.Service.InsertWord(word)

	words, err := s.Service.FindAllWords()
	if err != nil {
		// TODO Deal with error
	}
	var ws []core.Word
	for _, v := range words {
		ws = append(ws, *v)
	}
	components.WordForm(core.NewFormData()).Render(req.Context(), w)
	components.WordExtra(*word).Render(req.Context(), w)
}

func (s *handler) deleteWord(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
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

func (s *handler) getWordByDutch(w http.ResponseWriter, req *http.Request) {
	search := req.PathValue("text")
	word, err := s.Service.FindWordsBy(core.Search{DutchWord: &search})
	if err != nil {
		// TODO To improve error codes
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderJSON(w, word)
}

func (s *handler) getWords(w http.ResponseWriter, req *http.Request) {
	search := req.FormValue("word")
	words, err := s.Service.FindWordsBy(core.Search{DutchWord: &search})
	if err != nil {
		// TODO To improve error codes
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var ws []core.Word
	for _, v := range words {
		ws = append(ws, *v)
	}
	components.WordList(ws).Render(req.Context(), w)
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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
