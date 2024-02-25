package db

import (
	"errors"
	"fmt"
	"github.com/mikelangelon/dutch-words/core"
	"time"
)

var errNotFound = errors.New("not found")

type Store struct {
	memoryDB map[string]core.Word
}

func NewStore() Store {
	return Store{
		memoryDB: map[string]core.Word{},
	}
}
func (s Store) Insert(word *core.Word) error {
	word.ID = s.generateID()
	s.memoryDB[word.ID] = *word
	return nil
}

func (s Store) Delete(id string) error {
	delete(s.memoryDB, id)
	return nil
}

func (s Store) FindByID(id string) (*core.Word, error) {
	v, ok := s.memoryDB[id]
	if !ok {
		return nil, errNotFound
	}
	return &v, nil
}
func (s Store) FindByDutch(dutch string) (*core.Word, error) {
	for _, v := range s.memoryDB {
		if v.Dutch == dutch {
			return &v, nil
		}
	}
	return nil, errNotFound
}

func (s Store) FindAll() ([]*core.Word, error) {
	var words = make([]*core.Word, 0, len(s.memoryDB))
	for _, v := range s.memoryDB {
		words = append(words, &v)
	}
	return words, nil
}

func (s Store) generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixMicro())
}
