package server

import (
	"encoding/json"
	"github.com/mikelangelon/dutch-words/components"
	"github.com/mikelangelon/dutch-words/core"
	"github.com/mikelangelon/dutch-words/services"
	"net/http"
	"strings"
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
	tab1 := components.FormAndSearch(core.NewFormData(nil, s.getTags()), ws)
	components.Dashboard(navBar, tab1).Render(request.Context(), w)
}

func (s *handler) newTags(w http.ResponseWriter, request *http.Request) {
	enableCors(&w)
	tag := request.FormValue("tag")
	allTags := strings.Split(request.FormValue("all_tags"), ",")
	components.TagsField(core.FormData{
		Tags: append(allTags, tag),
	}).Render(request.Context(), w)
}
func (s *handler) getTags() []string {
	t, err := s.Service.FindAllTags()
	if err != nil {
		// TODO Deal with tags issue. Maybe skip it?
	}
	return t
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
	components.Dashboard(navBar, components.FormAndSearch(core.NewFormData(nil, s.getTags()), ws)).Render(request.Context(), w)
}

func (s *handler) createWord(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	dutch := req.FormValue("dutch")
	english := req.FormValue("english")
	tags := req.Form["tags"]
	wordType := req.FormValue("type")

	if dutch == "hond" {
		data := core.NewFormData(nil, s.getTags())
		data.Errors["word"] = "something crazy"
		components.WordForm(data).Render(req.Context(), w)
		http.Error(w, "duplicated", http.StatusUnprocessableEntity)
		return
	}
	word := core.NewWord(dutch, english, wordType, tags)
	err := s.Service.InsertWord(&word)

	words, err := s.Service.FindAllWords()
	if err != nil {
		// TODO Deal with error
	}
	var ws []core.Word
	for _, v := range words {
		ws = append(ws, *v)
	}
	components.WordForm(core.NewFormData(nil, s.getTags())).Render(req.Context(), w)
	components.WordExtra(word).Render(req.Context(), w)
}

func (s *handler) putWord(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	id := req.PathValue("id")
	dutch := req.FormValue("dutch")
	english := req.FormValue("english")
	tags := req.Form["tags"]

	word := &core.Word{
		ID:      id,
		Dutch:   dutch,
		English: english,
		Tags:    tags,
	}
	if err := s.Service.UpdateWord(word); err != nil {
		data := core.NewFormData(nil, s.getTags())
		data.Errors["word"] = "something crazy"
		components.WordForm(data).Render(req.Context(), w)
		http.Error(w, "duplicated", http.StatusUnprocessableEntity)
	}
	components.WordCard(*word).Render(req.Context(), w)
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
	components.WordCard(*word).Render(req.Context(), w)
}

func (s *handler) getWordEdit(w http.ResponseWriter, req *http.Request) {
	word, err := s.Service.FindWordByID(req.PathValue("id"))
	if err != nil {
		// TODO To improve error codes
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	components.WordCardEdit(core.NewFormData(word, s.getTags())).Render(req.Context(), w)
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

func (s *handler) renderTagsScreen(w http.ResponseWriter, req *http.Request) {
	tags, err := s.Service.FindAllTags()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	navBar := components.NavBar(nav("Tags"))
	components.Dashboard(navBar, components.Tags(tags)).Render(req.Context(), w)
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
