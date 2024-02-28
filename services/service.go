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
	FindByID(id string) (*core.Word, error)
	FindByDutch(dutch string) (*core.Word, error)
	FindAll() ([]*core.Word, error)
	FindBy(search core.Search) ([]*core.Word, error)
	Delete(id string) error
}

func (s Service) InsertWord(word *core.Word) error {
	return s.store.Insert(word)
}

func (s Service) DeleteWord(id string) error {
	return s.store.Delete(id)
}

func (s Service) FindWordByID(id string) (*core.Word, error) {
	return s.store.FindByID(id)
}

func (s Service) FindWordByDutch(dutch string) (*core.Word, error) {
	return s.store.FindByDutch(dutch)
}
func (s Service) FindAllWords() ([]*core.Word, error) {
	return s.store.FindAll()
}

func (s Service) FindWordsBy(search core.Search) ([]*core.Word, error) {
	return s.store.FindBy(search)
}
