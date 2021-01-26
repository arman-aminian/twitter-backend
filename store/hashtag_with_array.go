package store

// import (
// 	"container/heap"
// 	"github.com/arman-aminian/twitter-backend/model"
// 	"github.com/jinzhu/copier"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"time"
// )
//
// type HashtagStore struct {
// 	list        *[]*model.Hashtag
// 	lastUpdate  time.Time
// 	trends      *[]*model.Hashtag
// 	trendLength int
// }
//
// func NewHashtagStore() *HashtagStore {
// 	return &HashtagStore{
// 		list:        &[]*model.Hashtag{},
// 		lastUpdate:  time.Now(),
// 		trends:      &[]*model.Hashtag{},
// 		trendLength: 10,
// 	}
// }
//
// func (hh *HashtagStore) Len() int {
// 	return len(*hh.list)
// }
//
// func (hh *HashtagStore) Less(i, j int) bool {
// 	return (*hh.list)[i].Count > (*hh.list)[j].Count
// }
//
// func (hh *HashtagStore) Swap(i, j int) {
// 	(*hh.list)[i], (*hh.list)[j] = (*hh.list)[j], (*hh.list)[i]
// }
//
// func (hh *HashtagStore) Push(x interface{}) {
// 	*hh.list = append(*hh.list, x.(*model.Hashtag))
// }
//
// func (hh *HashtagStore) Pop() interface{} {
// 	old := *hh.list
// 	n := len(old)
// 	x := old[n-1]
// 	*hh.list = old[0 : n-1]
// 	return x
// }
//
// func (hh *HashtagStore) AddHashtag(ht *model.Hashtag) {
// 	exists := false
// 	for _, h := range *hh.list {
// 		if h.Name == ht.Name {
// 			// each pair of (tweet, hashtag) is unique so we don't need this.
// 			// for _, id := range *h.Tweets {
// 			// 	if id == (*ht.Tweets)[0] {
// 			// 		exists = true
// 			// 		break
// 			// 	}
// 			// }
// 			// if exists {
// 			// 	break
// 			// }
// 			ht.Count = h.Count + ht.Count
// 			*ht.Tweets = append(*h.Tweets, (*ht.Tweets)[0])
// 			exists = true
// 			break
// 		}
// 	}
// 	if exists {
// 		hh.RemoveHashtag(ht.Name)
// 		currTime := time.Now()
// 		diff := currTime.Sub(hh.lastUpdate)
// 		if diff.Hours() >= 1 {
// 			*hh.list = append(*hh.list, ht)
// 			hh.Update()
// 		} else {
// 			*hh.list = append(*hh.list, ht)
// 			if hh.Len() < hh.trendLength {
// 				hh.Update()
// 			} else if ht.Count > (*hh.trends)[len(*hh.trends)-1].Count {
// 				hh.Update()
// 			}
// 		}
// 	} else {
// 		*hh.list = append(*hh.list, ht)
// 		if hh.Len() < hh.trendLength {
// 			hh.Update()
// 		}
// 	}
// }
//
// func (hh *HashtagStore) GetHashtagByName(name string) *model.Hashtag {
// 	for _, h := range *hh.list {
// 		if h.Name == name {
// 			return h
// 		}
// 	}
// 	return nil
// }
//
// func (hh *HashtagStore) RemoveHashtag(name string) {
// 	newHeap := &[]*model.Hashtag{}
// 	for _, h := range *hh.list {
// 		if h.Name != name {
// 			*newHeap = append(*newHeap, h)
// 		}
// 	}
// 	*hh.list = *newHeap
// }
//
// func (hh *HashtagStore) DeleteTweetHashtags(t *model.Tweet, hashtags map[string]int) {
// 	for _, h := range *hh.list {
// 		if _, ok := hashtags[h.Name]; ok {
// 			newCount := h.Count
// 			newTweets := &[]primitive.ObjectID{}
// 			for _, id := range *h.Tweets {
// 				if id != t.ID {
// 					*newTweets = append(*newTweets, id)
// 				}
// 			}
// 			newCount -= hashtags[h.Name]
// 			h.Count = newCount
// 			h.Tweets = newTweets
// 		}
// 	}
// }
//
// func (hh *HashtagStore) Update() {
// 	hhCopy := NewHashtagStore()
// 	_ = copier.Copy(&hhCopy, &hh)
// 	heap.Init(hhCopy)
// 	for _, h := range *hh.list {
// 		heap.Push(hhCopy, h)
// 	}
//
// 	hh.trends = &[]*model.Hashtag{}
// 	for i := 0; i < hh.trendLength && hhCopy.Len() > 0; i++ {
// 		*hh.trends = append(*hh.trends, heap.Pop(hhCopy).(*model.Hashtag))
// 	}
//
// 	_ = copier.Copy(&hh, &hhCopy)
// 	hh.lastUpdate = time.Now()
// }
//
// func (hh *HashtagStore) GetTrends() *[]*model.Hashtag {
// 	return hh.trends
// }
