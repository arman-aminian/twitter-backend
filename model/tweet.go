package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tweet struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Text  string             `json:"text" bson:"text"`
	Media string             `json:"media" bson:"media"`
	Owner struct {
		Username       string `json:"username"`
		ProfilePicture string `json:"profile_picture"`
	} `json:"owner" bson:"owner"`
	Likes    *[]string `json:"likes" bson:"likes"`
	Retweets *[]string `json:"retweets" bson:"retweets"`
}

func NewTweet() *Tweet {
	var t Tweet
	t.Likes = &[]string{}
	t.Retweets = &[]string{}
	return &t
}
