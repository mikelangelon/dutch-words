package core

type Game struct {
	ID            string
	CurrentPoints int64
	Questions     []Question
	Next          bool
	Retry         bool
}

func (g Game) LatestQuestion() *Question {
	return &g.Questions[len(g.Questions)-1]
}
