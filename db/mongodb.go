package db

import (
	"context"
	"fmt"
	"github.com/mikelangelon/dutch-words/config"
	"github.com/mikelangelon/dutch-words/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
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
	Dutch      string   `bson:"dutch" json:"dutch"`
	English    string   `bson:"english" json:"english"`
	Difficulty *string  `bson:"difficulty,omitempty" json:"difficulty"`
	Tags       []string `bson:"tags,omitempty" json:"tags"`
}

func (m MongoStore) FindByID(id string) (*core.Word, error) {
	//TODO implement me
	panic("implement me")
}
func (m MongoStore) FindAll() ([]*core.Word, error) {
	return m.searchWords(nil)
}

func (m MongoStore) FindBy(search core.Search) ([]*core.Word, error) {
	var filter bson.M
	if search.Tag != nil {
		filter = bson.M{"tags": bson.M{"$in": []string{*search.Tag}}}
	} else if search.DutchWord != nil {
		filter = bson.M{"dutch": *search.DutchWord}
	} else if search.EnglishWord != nil {
		filter = bson.M{"english": *search.EnglishWord}
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
		words = append(words, &core.Word{
			ID:      w.Dutch,
			Dutch:   w.Dutch,
			English: w.English,
			Tags:    w.Tags,
		})
	}
	return words, nil
}

func (m MongoStore) Delete(id string) error {
	filter := bson.M{
		"id": id,
	}
	_, err := m.dutchCollection().DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("problem deleting task: %w", err)
	}
	return nil
}

func (m MongoStore) Insert(word *core.Word) error {
	ctx := context.TODO()
	if _, err := m.dutchCollection().InsertOne(ctx, word); err != nil {
		return fmt.Errorf("problem inserting word: %w", err)
	}
	return nil
}

func (m MongoStore) searchWords(limit *int64) ([]*core.Word, error) {
	var opts *options.FindOptions
	if limit != nil {
		opts = options.Find()
		opts.SetLimit(*limit)
	}
	c, err := m.dutchCollection().Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		return nil, fmt.Errorf("problem searching words: %w", err)
	}
	var words []*core.Word
	for c.Next(context.TODO()) {
		var w Word
		if err := c.Decode(&w); err != nil {
			return nil, err
		}
		words = append(words, &core.Word{
			ID:      w.Dutch,
			Dutch:   w.Dutch,
			English: w.English,
			Tags:    w.Tags,
		})
	}
	return words, nil
}

func (m MongoStore) dutchCollection() *mongo.Collection {
	return m.client.Database("dutch").Collection("words")
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
