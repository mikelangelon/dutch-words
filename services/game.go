package services

import (
	"github.com/mikelangelon/dutch-words/core"
	"log/slog"
)

type GameService struct {
	store store
}

func NewGameService(store store) GameService {
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
func (g GameService) loadPairsForTag(tag *string) ([]core.Word, string) {
	slog.Info("requesting for tag", "tag", tag)
	words, err := g.store.SearchBy(core.Search{Tag: tag})
	if err != nil {
		slog.Error("error getting words. loading hardcoded ones", "error", err)
		return []core.Word{}, ""
	}
	w := make([]core.Word, len(words))
	for i, v := range words {
		w[i] = core.Word{Dutch: v.Dutch, English: v.English}
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
		ID:            "A",
		Word:          words[0].English,
		CorrectOption: words[0].Dutch,
		Options:       options,
	}
}
