package store

import (
	"context"
	"github.com/arman-aminian/twitter-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (ts *TweetStore) GetTweetById(id *string) (*model.Tweet, error) {
	var t model.Tweet
	oid, err := primitive.ObjectIDFromHex(*id)
	if err != nil {
		return &t, nil
	}
	err = ts.db.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&t)
	return &t, err
}

func (ts *TweetStore) LikeTweet(t *model.Tweet, u *model.User) error {
	*t.Likes = append(*t.Likes, *model.NewOwner(u.Username, u.ProfilePicture))
	_, err := ts.db.UpdateOne(context.TODO(), bson.M{"_id": t.ID}, bson.M{"$set": bson.M{"likes": t.Likes}})
	if err != nil {
		return err
	}
	return nil
}

func (ts *TweetStore) UnLikeTweet(t *model.Tweet, u *model.User) error {
	newLikes := &[]model.Owner{}
	for _, o := range *t.Likes {
		if o.Username != u.Username {
			*newLikes = append(*newLikes, o)
		}
	}
	_, err := ts.db.UpdateOne(context.TODO(), bson.M{"_id": t.ID}, bson.M{"$set": bson.M{"likes": newLikes}})
	if err != nil {
		return err
	}
	t.Likes = newLikes
	return nil
}

func (ts *TweetStore) Retweet(t *model.Tweet, u *model.User) error {
	*t.Retweets = append(*t.Retweets, *model.NewOwner(u.Username, u.ProfilePicture))
	_, err := ts.db.UpdateOne(context.TODO(), bson.M{"_id": t.ID}, bson.M{"$set": bson.M{"retweets": t.Retweets}})
	if err != nil {
		return err
	}
	return nil
}

func (ts *TweetStore) UnRetweet(t *model.Tweet, u *model.User) error {
	newRetweets := &[]model.Owner{}
	for _, o := range *t.Retweets {
		if o.Username != u.Username {
			*newRetweets = append(*newRetweets, o)
		}
	}
	_, err := ts.db.UpdateOne(context.TODO(), bson.M{"_id": t.ID}, bson.M{"$set": bson.M{"retweets": newRetweets}})
	if err != nil {
		return err
	}
	t.Retweets = newRetweets
	return nil
}
