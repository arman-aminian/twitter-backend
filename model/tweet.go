package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Tweet struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Text     string             `json:"text" bson:"text"`
	Media    string             `json:"media" bson:"media"`
	Time     time.Time          `json:"time" bson:"time"`
	Owner    Owner              `json:"owner" bson:"owner"`
	Likes    *[]Owner           `json:"likes" bson:"likes"`
	Retweets *[]Owner           `json:"retweets" bson:"retweets"`
}

func NewTweet() *Tweet {
	var t Tweet
	t.Likes = &[]Owner{}
	t.Retweets = &[]Owner{}
	return &t
}
