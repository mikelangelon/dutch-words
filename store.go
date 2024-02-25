package main

import (
	"errors"
	"fmt"
	"time"
)

var errNotFound = errors.New("not found")

type store struct {
	memoryDB map[string]Word
}

func NewStore() store {
	return store{
		memoryDB: map[string]Word{},
	}
}
func (s store) Insert(word Word) error {
	word.ID = s.generateID()
	s.memoryDB[word.ID] = word
	return nil
}

func (s store) Delete(id string) error {
	delete(s.memoryDB, id)
	return nil
}

func (s store) FindByID(id string) (*Word, error) {
	v, ok := s.memoryDB[id]
	if !ok {
		return nil, errNotFound
	}
	return &v, nil
}
func (s store) FindByDutch(dutch string) (*Word, error) {
	for _, v := range s.memoryDB {
		if v.Dutch == dutch {
			return &v, nil
		}
	}
	return nil, errNotFound
}

func (s store) FindAll() ([]*Word, error) {
	var words = make([]*Word, 0, len(s.memoryDB))
	for _, v := range s.memoryDB {
		words = append(words, &v)
	}
	return words, nil
}

func (s store) generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixMicro())
}
