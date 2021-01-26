package hashtag

import (
	"github.com/arman-aminian/twitter-backend/model"
)

type Store interface {
	// Max heap is shit I think. Going with classic linear array
	AddHashtag(ht *model.Hashtag) error                   // O(n + O(exists))
	GetHashtagByName(name string) (*model.Hashtag, error) // O(n)
	RemoveHashtag(name string) error                      // O(n)
	DeleteTweetHashtags(t *model.Tweet, hashtags map[string]int) error
	// DeleteOldHashtags() // Should be in handler
	Update() error
	GetTrends() *[]*model.Hashtag
}
