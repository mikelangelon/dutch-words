package services

import (
	"log/slog"

	"github.com/mikelangelon/dutch-words/core"
)

type GameService struct {
	store gameStore
}

type gameStore interface {
	SearchBy(search core.Search) ([]*core.Word, error)
	UpsertAnswer(a core.Answer) error
}

func NewGameService(store gameStore) GameService {
	return GameService{store: store}
}

func (g GameService) NewGame() core.Game {
	return core.Game{
		ID: "ABCD",
		Questions: []core.Question{
			g.NextQuestion(),
		},
	}
}

func (g GameService) NextQuestion() core.Question {
	words, _ := g.loadPairsForTag(nil)
	return ToQuestion(words)
}

func (g GameService) Answer(answer core.Answer) error {
	return g.store.UpsertAnswer(answer)
}
func (g GameService) loadPairsForTag(tag *string) ([]core.Word, string) {
	slog.Info("requesting for tag", "tag", tag)
	words, err := g.store.SearchBy(core.Search{Tag: tag})
	if err != nil {
		slog.Error("error getting words. loading hardcoded ones", "error", err)
		return []core.Word{}, ""
	}
	w := make([]core.Word, len(words))
	for i, v := range words {
		w[i] = core.Word{ID: v.ID, Dutch: v.Dutch, English: v.English}
	}
	return w, ""
}

func ToQuestion(words []core.Word) core.Question {
	options := []core.Option{
		{Text: words[0].Dutch},
		{Text: words[1].Dutch},
		{Text: words[2].Dutch},
		{Text: words[3].Dutch},
	}
	core.ShuffleOption(options)
	return core.Question{
		WordID:         words[0].ID,
		Word:           words[0].English,
		CorrectOption:  words[0].Dutch,
		Options:        options,
		QuestionPoints: 1,
	}
}
