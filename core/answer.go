package core

type Answer struct {
	WordID  string `json:"wordId"`
	UserID  string `json:"userId"`
	Correct bool   `json:"correct"`
}
