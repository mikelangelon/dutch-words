package core

type SentenceData struct {
	Sentence Sentence
}

func NewSentenceData(sentence Sentence) SentenceData {
	return SentenceData{
		Sentence: sentence,
	}
}
