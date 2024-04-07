package core

import "math/rand"

type Question struct {
	ID            string
	Word          string
	CorrectOption string
	Options       []Option
	Next          bool
	Retry         bool
}

type Option struct {
	Text   string
	Status int64
}

func ShuffleOption(options []Option) {
	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})
}
