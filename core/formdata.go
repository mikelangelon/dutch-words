package core

type FormData struct {
	Word   *Word
	Tags   []string
	Errors map[string]string
}

func NewFormData(word *Word, tags []string) FormData {
	return FormData{
		Tags:   tags,
		Word:   word,
		Errors: make(map[string]string),
	}
}
