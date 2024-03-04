package core

import (
	"fmt"
	"slices"
	"time"
)

type Word struct {
	ID      string
	Dutch   string
	English string
	Type    string
	Tags    []string
}

func NewWord(dutch, english, wordType string, tags []string) Word {
	return Word{
		ID:      fmt.Sprintf("%d", time.Now().UnixNano()),
		Dutch:   dutch,
		English: english,
		Type:    wordType,
		Tags:    tags,
	}
}

func (w *Word) HasTag(tag string) bool {
	if w == nil {
		return false
	}
	return slices.Contains(w.Tags, tag)
}
