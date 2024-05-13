package services

import "github.com/mikelangelon/dutch-words/core"

type scoreStore interface {
	GetAnswers() ([]core.Answer, error)
}
type ScoreService struct {
	store scoreStore
}

func NewScoreService(store scoreStore) ScoreService {
	return ScoreService{store: store}
}
