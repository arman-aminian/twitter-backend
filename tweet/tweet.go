package tweet

import (
	"github.com/arman-aminian/twitter-backend/model"
)

type Store interface {
	CreateTweet(*model.Tweet) error
	GetTweetById(id *string) (*model.Tweet, error)
	LikeTweet(t *model.Tweet, u *model.User) error
	UnLikeTweet(t *model.Tweet, u *model.User) error
	Retweet(t *model.Tweet, u *model.User) error
	UnRetweet(t *model.Tweet, u *model.User) error
}
