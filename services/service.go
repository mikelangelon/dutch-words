package services

import (
	"github.com/mikelangelon/dutch-words/core"
	"github.com/mikelangelon/dutch-words/db"
)

type Service struct {
	store db.Store
}

func NewService(store db.Store) Service {
	return Service{store: store}
}

func (s Service) InsertWord(word core.Word) error {
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
