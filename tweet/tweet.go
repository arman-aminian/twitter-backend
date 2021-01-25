package tweet

import (
	"github.com/arman-aminian/twitter-backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	CreateTweet(*model.Tweet) error
	RemoveTweet(*model.Tweet) error
	GetTweetById(id *string) (*model.Tweet, error)
	GetAllTweets() ([]bson.M, error)
	LikeTweet(t *model.Tweet, u *model.User) error
	UnLikeTweet(t *model.Tweet, u *model.User) error
	Retweet(t *model.Tweet, u *model.User) error
	UnRetweet(t *model.Tweet, u *model.User) error

	ExtractHashtags(t *model.Tweet) map[string]int

	GetTimelineFromFollowingsUsernames(usernames []string) (*[]model.Tweet, error)
}
