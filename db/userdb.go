package db

import (
	"github.com/mikelangelon/dutch-words/core"
	"go.mongodb.org/mongo-driver/bson"
)

type Answer struct {
	WordID  string `bson:"wordId" json:"wordId"`
	UserID  string `bson:"userId" json:"userId"`
	Correct int64  `bson:"amountCorrect" json:"amountCorrect"`
	Wrong   int64  `bson:"amountWrong" json:"amountWrong"`
}

func increment(a core.Answer) bson.D {
	if a.Correct {
		return bson.D{{Key: "amountCorrect", Value: 1}}
	}
	return bson.D{{Key: "amountWrong", Value: 1}}
}
