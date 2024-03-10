package services

import "github.com/mikelangelon/dutch-words/core"

type sentencesStore interface {
	InsertSentence(s *core.Sentence) error
	UpdateSentence(s *core.Sentence) error
	FindSentencesBy() ([]*core.Sentence, error)
	DeleteSentence(id string) error
}

type SentencesService struct {
	store sentencesStore
}

func NewSentencesService(store sentencesStore) SentencesService {
	return SentencesService{store: store}
}

func (s SentencesService) Insert(sentence *core.Sentence) error {
	return s.store.InsertSentence(sentence)
}

func (s SentencesService) Update(sentence *core.Sentence) error {
	return s.store.UpdateSentence(sentence)
}

func (s SentencesService) Delete(id string) error {
	return s.store.DeleteSentence(id)
}
func (s SentencesService) FindAll() ([]*core.Sentence, error) {
	return s.store.FindSentencesBy()
}
