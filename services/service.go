package services

import (
	"github.com/mikelangelon/dutch-words/core"
)

type Service struct {
	store store
}

func NewService(store store) Service {
	return Service{store: store}
}

type store interface {
	Insert(word *core.Word) error
	FindBy(search core.Search) ([]*core.Word, error)
	Delete(id string) error
	GetAllTags() ([]string, error)
}

func (s Service) InsertWord(word *core.Word) error {
	return s.store.Insert(word)
}

func (s Service) DeleteWord(id string) error {
	return s.store.Delete(id)
}

func (s Service) FindWordByID(id string) (*core.Word, error) {
	ws, err := s.store.FindBy(core.Search{ID: &id})
	if err != nil {
		return nil, err
	}
	return ws[0], nil
}

func (s Service) FindAllWords() ([]*core.Word, error) {
	return s.store.FindBy(core.Search{})
}

func (s Service) FindWordsBy(search core.Search) ([]*core.Word, error) {
	return s.store.FindBy(search)
}

func (s Service) FindAllTags() ([]string, error) {
	return s.store.GetAllTags()
}
