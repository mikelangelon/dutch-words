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
	Types   []string
	Article *string
	Tags    []string
}

func NewWord(dutch, english string, wordType, tags []string) Word {
	return Word{
		ID:      fmt.Sprintf("%d", time.Now().UnixNano()),
		Dutch:   dutch,
		English: english,
		Types:   wordType,
		Tags:    tags,
	}
}

func (w *Word) HasTag(tag string) bool {
	return w.has(tag, w.Tags)
}

func (w *Word) HasType(s string) bool {
	if w == nil {
		return false
	}
	return w.has(s, w.Types)
}

func (w *Word) has(s string, values []string) bool {
	if w == nil {
		return false
	}
	return slices.Contains(values, s)
}
