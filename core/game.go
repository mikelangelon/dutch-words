package core

type Game struct {
	ID        string
	Questions []Question
}

func (g Game) LatestQuestion() Question {
	return g.Questions[len(g.Questions)-1]
}
