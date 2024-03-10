package db

import (
	"context"
	"fmt"
	"github.com/mikelangelon/dutch-words/core"
	"go.mongodb.org/mongo-driver/bson"
)

type Sentence struct {
	ID      string `bson:"_id" json:"id"`
	Dutch   string `bson:"dutch" json:"dutch"`
	English string `bson:"english" json:"english"`
}

func (w Sentence) toEntity() *core.Sentence {
	return &core.Sentence{
		ID:      w.ID,
		Dutch:   w.Dutch,
		English: w.English,
	}
}

func sentenceToDB(w *core.Sentence) Sentence {
	return Sentence{
		ID:      w.ID,
		Dutch:   w.Dutch,
		English: w.English,
	}
}

func (m MongoStore) InsertSentence(s *core.Sentence) error {
	ctx := context.TODO()
	if _, err := m.sentencesCollection().InsertOne(ctx, sentenceToDB(s)); err != nil {
		return fmt.Errorf("problem inserting sentence: %w", err)
	}
	return nil
}

func (m MongoStore) UpdateSentence(s *core.Sentence) error {
	ctx := context.TODO()
	if _, err := m.sentencesCollection().UpdateByID(ctx, s.ID, bson.M{"$set": sentenceToDB(s)}); err != nil {
		return fmt.Errorf("problem inserting sentence: %w", err)
	}
	return nil
}
func (m MongoStore) DeleteSentence(id string) error {
	filter := bson.M{
		"_id": id,
	}
	_, err := m.sentencesCollection().DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("problem deleting sentence: %w", err)
	}
	return nil
}

func (m MongoStore) FindSentencesBy() ([]*core.Sentence, error) {
	var filter bson.M
	c, err := m.sentencesCollection().Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("problem searching sentence: %w", err)
	}
	var sentences []*core.Sentence
	for c.Next(context.TODO()) {
		var s Sentence
		if err := c.Decode(&s); err != nil {
			return nil, err
		}
		sentences = append(sentences, s.toEntity())
	}
	return sentences, nil
}
