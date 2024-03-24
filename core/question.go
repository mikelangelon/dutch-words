package core

type Question struct {
	ID      string
	Word    string
	Options []Option
	Next    bool
	Retry   bool
}

type Option struct {
	Text   string
	Status int64
}
