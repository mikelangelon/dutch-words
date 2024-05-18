package core

type Answer struct {
	WordID  string `json:"wordId"`
	UserID  string `json:"userId"`
	Correct bool   `json:"correct"`

	Word          Word `json:"word"`
	AmountCorrect int  `json:"amountCorrect"`
	AmountWrong   int  `json:"amountWrong"`
}
