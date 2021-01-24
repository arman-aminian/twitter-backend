package store

import (
	"context"
	"github.com/arman-aminian/twitter-backend/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type TweetStore struct {
	db *mongo.Collection
}

func NewTweetStore(db *mongo.Collection) *TweetStore {
	return &TweetStore{
		db: db,
	}
}

func (ts *TweetStore) CreateTweet(t *model.Tweet) error {
	_, err := ts.db.InsertOne(context.TODO(), t)
	return err
}
