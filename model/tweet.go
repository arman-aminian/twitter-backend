package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Tweet struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Text     string             `json:"text" bson:"text"`
	Media    string             `json:"media" bson:"media"`
	Date     string             `json:"date" bson:"date"`
	Time     time.Time          `json:"time" bson:"time"`
	Owner    Owner              `json:"owner" bson:"owner"`
	Likes    *[]Owner           `json:"likes" bson:"likes"`
	Retweets *[]Owner           `json:"retweets" bson:"retweets"`
	Parents  *[]CommentTweet    `json:"parents" bson:"parents"`
	Comments *[]CommentTweet    `json:"comments" bson:"comments"`
}

type TweetList struct {
	Tweets *[]Tweet `json:"tweets"`
}

type TweetIdList struct {
	Tweets []string `json:"tweets"`
}

type CommentTweet struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Text          string             `json:"text" bson:"text"`
	Media         string             `json:"media" bson:"media"`
	Date          string             `json:"date" bson:"date"`
	Time          time.Time          `json:"time" bson:"time"`
	Owner         Owner              `json:"owner" bson:"owner"`
	Likes         *[]Owner           `json:"likes" bson:"likes"`
	Retweets      *[]Owner           `json:"retweets" bson:"retweets"`
	CommentsCount int                `json:"comments_count"`
}

func NewTweet() *Tweet {
	var t Tweet
	t.Likes = &[]Owner{}
	t.Retweets = &[]Owner{}
	t.Parents = &[]CommentTweet{}
	t.Comments = &[]CommentTweet{}
	return &t
}

func NewCommentTweet(tweet Tweet) *CommentTweet {
	var c CommentTweet
	c.ID = tweet.ID
	c.Text = tweet.Text
	c.Media = tweet.Media
	c.Date = tweet.Date
	c.Time = tweet.Time
	c.Owner = tweet.Owner
	c.CommentsCount = len(*tweet.Comments)
	return &c
}
