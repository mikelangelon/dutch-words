package db

import (
	"context"
	"fmt"
	"github.com/mikelangelon/dutch-words/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type answerDTO struct {
}

func (m MongoStore) answersCollection() *mongo.Collection {
	return m.client.Database("dutch").Collection("answer")
}

func (m MongoStore) UpsertAnswer(a core.Answer) error {
	ctx := context.TODO()
	objectId, err := primitive.ObjectIDFromHex(a.WordID)
	if err != nil {
		return err
	}
	_, err = m.answersCollection().UpdateOne(ctx, bson.M{"wordId": objectId, "userId": a.UserID},
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
	findOptions.SetSort(bson.D{{"amountWrong", -1}})
	//cursor, err := m.answersCollection().Find(context.TODO(), bson.M{"amountCorrect": bson.M{"$gt": 0}}, findOptions)
	//if err != nil {
	//	return nil, err
	//}

	pipeline := []bson.M{
		bson.M{"$lookup": bson.M{"from": "words", "localField": "wordId", "foreignField": "_id", "as": "words"}},
		{"$set": bson.M{"word": bson.M{"$arrayElemAt": bson.A{"$words", 0}}}},
		//{"$set": bson.M{"dutch": "$word.dutch"}}, //word:  { $arrayElemAt: [ "$words", 0 ] }},
		{"$sort": bson.M{"amountWrong": -1}},
	}
	cursor, err := m.answersCollection().Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	var answers []core.Answer
	if err = cursor.All(context.TODO(), &answers); err != nil {
		panic(err)
	}
	//var answers []Answer
	//if err = cursor.All(context.TODO(), &answers); err != nil {
	//	panic(err)
	//}
	return answers, nil
}

//func (m MongoStore) GetAnswers() ([]core.Answer, error) {
//	findOptions := options.Find()
//	findOptions.SetSort(bson.D{{"amountCorrect", -1}})
//	cursor, err := m.answersCollection().Find(context.TODO(), bson.M{"amountCorrect": bson.M{"$gt": 0}}, findOptions)
//	if err != nil {
//		return nil, err
//	}
//	defer cursor.Close(context.TODO())
//	var answers []Answer
//	if err = cursor.All(context.TODO(), &answers); err != nil {
//		panic(err)
//	}
//
//	var result []core.Answer
//	for _, v := range answers {
//		words, err := m.FindBy(core.Search{ID: &v.UserID})
//		if err != nil {
//			continue
//		}
//		if len(words) == 0 {
//			continue
//		}
//		result = append(result, core.Answer{
//			WordID:  "",
//			Correct: false,
//			Word:    words[0].Dutch,
//		})
//
//	}
//	return result, nil
//}
