package db

import (
	"context"
	"fmt"

	"github.com/mikelangelon/dutch-words/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Answer struct {
	WordID  string `bson:"wordId" json:"wordId"`
	UserID  string `bson:"userId" json:"userId"`
	Correct int64  `bson:"amountCorrect" json:"amountCorrect"`
	Wrong   int64  `bson:"amountWrong" json:"amountWrong"`
}

func (m MongoStore) UpsertAnswer(a core.Answer) error {
	ctx := context.TODO()

	_, err := m.answersCollection().UpdateOne(ctx, bson.M{"wordId": a.WordID, "userId": a.UserID},
		bson.D{
			{Key: "$inc", Value: increment(a)},
		}, options.Update().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("problem upserting answer: %w", err)
	}
	return nil
}

func increment(a core.Answer) bson.D {
	if a.Correct {
		return bson.D{{Key: "amountCorrect", Value: 1}}
	}
	return bson.D{{Key: "amountWrong", Value: 1}}
}
