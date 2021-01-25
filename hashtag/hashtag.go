package hashtag

import (
	"github.com/arman-aminian/twitter-backend/model"
)

type Store interface {
	// Max heap is shit I think. Going with classic linear array
	AddHashtag(ht *model.Hashtag)                // O(n + O(exists))
	GetHashtagByName(name string) *model.Hashtag // O(n)
	RemoveHashtag(name string)                   // O(n)
	DeleteTweetHashtags(t *model.Tweet, hashtags map[string]int)
	// DeleteOldHashtags() // Should be in handler
	Update()
	GetTrends() *[]*model.Hashtag
}
