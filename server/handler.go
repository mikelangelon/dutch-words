package server

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/a-h/templ"

	"github.com/mikelangelon/dutch-words/components"
	"github.com/mikelangelon/dutch-words/core"
	"github.com/mikelangelon/dutch-words/services"
)

type handler struct {
	Service          services.Service
	SentencesService services.SentencesService
}

func newHandler(service services.Service, sentencesService services.SentencesService) *handler {
	return &handler{
		Service:          service,
		SentencesService: sentencesService,
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
	err = components.WordList(ws).Render(request.Context(), w)
	if err != nil {
		slog.Error("problem rendering", "error", err)
	}
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
	err := components.Dashboard(navBar, tab1).Render(request.Context(), w)
	if err != nil {
		slog.Error("problem rendering", "error", err)
	}
}

func (s *handler) newTags(w http.ResponseWriter, request *http.Request) {
	enableCors(&w)
	tag := request.FormValue("tag")
	allTags := strings.Split(request.FormValue("all_tags"), ",")
	err := components.TagsField(core.FormData{
		Tags: append(allTags, tag),
	}).Render(request.Context(), w)
	if err != nil {
		slog.Error("problem rendering", "error", err)
	}
}
func (s *handler) getTags() []string {
	t, err := s.Service.FindAllTags()
	if err != nil {
		// TODO Deal with tags issue. Maybe skip it?
		slog.Error("problem getting tags", "error", err)
	}
	return t
}

func nav(current string) core.NavigationItems {
	items := core.NavigationItems{
		{
			Label: "Home",
			Link:  "/web/tab1",
		},
		{
			Label: "Sentences",
			Link:  "/web/sentence-tab",
		},
		{
			Label: "Search",
			Link:  "/web/tab2",
		},
		{
			Label: "Tags",
			Link:  "/web/tab3",
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
	err := components.Dashboard(navBar, components.FormAndSearch(core.NewFormData(nil, s.getTags()), ws)).Render(request.Context(), w)
	if err != nil {
		slog.Error("problem rendering", "error", err)
	}
}
func (s *handler) createWord(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	dutch := req.FormValue("dutch")
	english := req.FormValue("english")
	tags := req.Form["tags"]
	types := req.Form["types"]

	if dutch == "hond" {
		data := core.NewFormData(nil, s.getTags())
		data.Errors["word"] = "something crazy"
		err := components.WordForm(data).Render(req.Context(), w)
		if err != nil {
			slog.Error("problem rendering", "error", err)
		}
		http.Error(w, "duplicated", http.StatusUnprocessableEntity)
		return
	}
	word := core.NewWord(dutch, english, types, tags)
	err := s.Service.InsertWord(&word)
	if err != nil {
		slog.Error("problem inserting word", "error", err)
		http.Error(w, "duplicated", http.StatusInternalServerError)
		return
	}
	err = components.WordForm(core.NewFormData(nil, s.getTags())).Render(req.Context(), w)
	if err != nil {
		slog.Error("problem rendering word form", "error", err)
	}
	err = components.WordExtra(word).Render(req.Context(), w)
	if err != nil {
		slog.Error("problem rendering extra word", "error", err)
	}
}

func (s *handler) putWord(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	id := req.PathValue("id")
	dutch := req.FormValue("dutch")
	english := req.FormValue("english")
	types := req.Form["types"]
	tags := req.Form["tags"]

	word := &core.Word{
		ID:      id,
		Dutch:   dutch,
		English: english,
		Types:   types,
		Tags:    tags,
	}
	if err := s.Service.UpdateWord(word); err != nil {
		data := core.NewFormData(nil, s.getTags())
		data.Errors["word"] = "something crazy"
		err := components.WordForm(data).Render(req.Context(), w)
		if err != nil {
			slog.Error("problem rendering", "error", err)
		}
		http.Error(w, "duplicated", http.StatusUnprocessableEntity)
	}
	err := components.WordCard(*word).Render(req.Context(), w)
	if err != nil {
		slog.Error("problem rendering word card", "error", err)
	}
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
	err = components.WordCard(*word).Render(req.Context(), w)
	if err != nil {
		return
	}
}

func (s *handler) getWordEdit(w http.ResponseWriter, req *http.Request) {
	word, err := s.Service.FindWordByID(req.PathValue("id"))
	if err != nil {
		// TODO To improve error codes
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = components.WordCardEdit(core.NewFormData(word, s.getTags())).Render(req.Context(), w)
	if err != nil {
		slog.Error("problem rendering word card edit")
	}
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
	err = components.WordList(ws).Render(req.Context(), w)
	if err != nil {
		slog.Error("problem rendering word list", "error", err)
	}
}

func (s *handler) renderTagsScreen(w http.ResponseWriter, req *http.Request) {
	tags, err := s.Service.FindAllTags()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	navBar := components.NavBar(nav("Tags"))
	err = components.Dashboard(navBar, components.Tags(tags)).Render(req.Context(), w)
	if err != nil {
		slog.Error("problem rendering dashboard", "error", err)
	}
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		slog.Error("problem rendering JSON", "error", err)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func formAndList(sentences []*core.Sentence) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		err := components.SentenceEdit(core.SentenceData{Sentence: core.Sentence{}}).Render(context.TODO(), w)
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, "<div>")
		if err != nil {
			return err
		}
		err = components.SentencesList(sentences).Render(context.TODO(), w)
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, "</div>")
		if err != nil {
			return err
		}
		return nil
	})
}
