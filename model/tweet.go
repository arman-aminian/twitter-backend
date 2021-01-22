package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tweet struct {
	Text     string                `json:"text" bson:"text"`
	Media    string                `json:"media" bson:"media"`
	Owner    string                `json:"owner" bson:"owner"`
	Likes    *[]primitive.ObjectID `json:"likes" bson:"likes"`
	Retweets *[]primitive.ObjectID `json:"retweets" bson:"retweets"`
}
