package tweet

import (
	"github.com/arman-aminian/twitter-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Store interface {
	CreateTweet(*model.Tweet) error
	AddCommentToTweet(parent *model.Tweet, child *model.CommentTweet) error
	RemoveTweet(*model.Tweet) error
	GetTweetById(id *string) (*model.Tweet, error)
	GetAllTweets() ([]bson.M, error)
	LikeTweet(t *model.Tweet, u *model.User) error
	UnLikeTweet(t *model.Tweet, u *model.User) error
	Retweet(t *model.Tweet, u *model.User) error
	UnRetweet(t *model.Tweet, u *model.User) error

	ExtractHashtags(t *model.Tweet) map[string]int

	GetTimelineFromTweetIDs(usernames []primitive.ObjectID, day int) (*[]model.Tweet, error)
	GetTweetSearchResult(username string) (*[]model.Tweet, error)
}
