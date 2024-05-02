package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mikelangelon/dutch-words/config"
	"github.com/mikelangelon/dutch-words/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoStore struct {
	client *mongo.Client
}

func NewMongoStore(cfg *config.Config) (MongoStore, error) {
	client, err := setup(cfg.MongoURL)
	if err != nil {
		return MongoStore{}, err
	}
	return MongoStore{client: client}, nil
}

type Word struct {
	ID         string   `bson:"_id" json:"id"`
	Dutch      string   `bson:"dutch" json:"dutch"`
	English    string   `bson:"english" json:"english"`
	Difficulty *string  `bson:"difficulty,omitempty" json:"difficulty"`
	Types      []string `bson:"types,omitempty" json:"types"`
	Tags       []string `bson:"tags,omitempty" json:"tags"`
	Article    *string  `bson:"article,omitempty" json:"article"`
}

func (w Word) toEntity() *core.Word {
	return &core.Word{
		ID:      w.ID,
		Dutch:   w.Dutch,
		English: w.English,
		Types:   w.Types,
		Tags:    w.Tags,
		Article: w.Article,
	}
}

func wordToDB(w *core.Word) Word {
	return Word{
		ID:      w.ID,
		Dutch:   w.Dutch,
		English: w.English,
		Types:   w.Types,
		Tags:    w.Tags,
		Article: w.Article,
	}
}
func (m MongoStore) Insert(word *core.Word) error {
	ctx := context.TODO()
	if _, err := m.dutchCollection().InsertOne(ctx, wordToDB(word)); err != nil {
		return fmt.Errorf("problem inserting word: %w", err)
	}
	return nil
}

func (m MongoStore) Update(word *core.Word) error {
	ctx := context.TODO()
	if _, err := m.dutchCollection().UpdateByID(ctx, word.ID, bson.M{"$set": wordToDB(word)}); err != nil {
		return fmt.Errorf("problem inserting word: %w", err)
	}
	return nil
}

func (m MongoStore) Delete(id string) error {
	filter := bson.M{
		"_id": id,
	}
	_, err := m.dutchCollection().DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("problem deleting task: %w", err)
	}
	return nil
}

type countResults struct {
	Tag   string `bson:"_id"`
	Count int    `bson:"count"`
}

func (m MongoStore) GetAllTags() (core.Tags, error) {

	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$tags"},
			{"count", bson.D{{"$sum", 1}}},
		}}}

	cursor, err := m.dutchCollection().Aggregate(context.TODO(), mongo.Pipeline{bson.D{{"$unwind", "$tags"}}, groupStage})
	if err != nil {
		return nil, err
	}

	// display the results
	var results []countResults
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	var count []core.CountTag
	for _, r := range results {
		count = append(count, core.CountTag{Tag: r.Tag, Count: r.Count})
	}
	return count, nil
}

func (m MongoStore) FindBy(search core.Search) ([]*core.Word, error) {
	var filter bson.M
	if search.Tag != nil {
		filter = bson.M{"tags": bson.M{"$in": []string{*search.Tag}}}
	} else if search.DutchWord != nil {
		filter = bson.M{"dutch": bson.D{{Key: "$regex", Value: fmt.Sprintf("^%s", *search.DutchWord)}}}
	} else if search.EnglishWord != nil {
		filter = bson.M{"english": bson.D{{Key: "$regex", Value: fmt.Sprintf("^%s", *search.EnglishWord)}}}
	} else if search.ID != nil {
		filter = bson.M{"_id": *search.ID}
	}

	c, err := m.dutchCollection().Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("problem searching words: %w", err)
	}
	var words []*core.Word
	for c.Next(context.TODO()) {
		var w Word
		if err := c.Decode(&w); err != nil {
			return nil, err
		}
		words = append(words, w.toEntity())
	}
	return words, nil
}

// TODO Almost the same as FindBy, decide which one to keep
func (m MongoStore) SearchBy(search core.Search) ([]*core.Word, error) {
	if search.Limit == 0 {
		search.Limit = 20
	}
	sample := bson.D{{Key: "$sample", Value: bson.D{
		{Key: "size", Value: search.Limit},
	}}}
	var s = bson.A{}
	if search.Tag != nil {
		s = append(s, bson.D{{Key: "$match", Value: bson.D{{Key: "tags", Value: *search.Tag}}}})
	}
	s = append(s, sample)
	ctx := context.TODO()
	c, err := m.dutchCollection().Aggregate(ctx, s)
	if err != nil {
		return nil, err
	}
	defer func(c *mongo.Cursor, ctx context.Context) {
		closeErr := c.Close(ctx)
		if closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}(c, ctx)

	var words []*core.Word
	for c.Next(context.TODO()) {
		var w Word
		if err := c.Decode(&w); err != nil {
			return nil, err
		}
		words = append(words, w.toEntity())
	}
	return words, nil
}
func (m MongoStore) dutchCollection() *mongo.Collection {
	return m.client.Database("dutch").Collection("words")
}

func (m MongoStore) sentencesCollection() *mongo.Collection {
	return m.client.Database("dutch").Collection("sentences")
}

func (m MongoStore) answersCollection() *mongo.Collection {
	return m.client.Database("dutch").Collection("answer")
}

func setup(uri string) (*mongo.Client, error) {
	clientOptions := options.Client().
		ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	if err := c.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}
	return c, nil
}
