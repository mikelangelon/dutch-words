package db

import (
	"context"
	"errors"
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

func (m MongoStore) SearchWords() ([]*core.Word, error) {
	options := options.Find()
	options.SetLimit(10)
	c, err := m.client.Database("dutch").Collection("words").Find(context.TODO(), bson.D{}, options)
	if err != nil {
		return nil, errors.New("crash")
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
