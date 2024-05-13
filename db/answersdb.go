package db

import (
	"context"
	"fmt"
	"github.com/mikelangelon/dutch-words/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m MongoStore) answersCollection() *mongo.Collection {
	return m.client.Database("dutch").Collection("answer")
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

func (m MongoStore) GetAnswers() ([]core.Answer, error) {
	findOptions := options.Find()
	// Sort by `price` field descending
	findOptions.SetSort(bson.D{{"amountCorrect", -1}})
	cursor, err := m.answersCollection().Find(context.TODO(), bson.M{"amountCorrect": bson.M{"$gt": 0}}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var answers []Answer
	if err = cursor.All(context.TODO(), &answers); err != nil {
		panic(err)
	}
	return nil, nil
}
