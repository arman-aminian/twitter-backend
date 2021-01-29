package hashtag

import (
	"github.com/arman-aminian/twitter-backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Store interface {
	AddHashtag(ht *model.Hashtag) error                   // O(n + O(exists))
	GetHashtagByName(name string) (*model.Hashtag, error) // O(n)
	RemoveHashtag(name string) error                      // O(n)
	DeleteTweetHashtags(t *model.Tweet, hashtags map[string]int) error
	GetHashtagTweets(name string) (*[]primitive.ObjectID, error)
	Update() error
	GetTrends() *[]*model.Hashtag
}
