package store

import (
	"context"
	"fmt"
	"github.com/arman-aminian/twitter-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type HashtagStore struct {
	db          *mongo.Collection
	lastUpdate  time.Time
	trends      *[]*model.Hashtag
	trendLength int64
}

func NewHashtagStore(db *mongo.Collection) *HashtagStore {
	return &HashtagStore{
		db:          db,
		lastUpdate:  time.Now(),
		trends:      &[]*model.Hashtag{},
		trendLength: 10,
	}
}

func (hs *HashtagStore) AddHashtag(ht *model.Hashtag) error {
	exists := false
	h, err := hs.GetHashtagByName(ht.Name)
	if err == mongo.ErrNoDocuments {
		_, err := hs.db.InsertOne(context.TODO(), ht)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		exists = true
		newCount := h.Count + ht.Count
		newTweets := append(*h.Tweets, *ht.Tweets...)
		_, err = hs.db.UpdateOne(context.TODO(), bson.M{"name": h.Name}, bson.M{"$set": bson.M{"count": newCount, "tweets": newTweets}})
		if err != nil {
			return err
		}
	}
	currLen, err := hs.db.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return err
	}
	if exists {
		currTime := time.Now()
		diff := currTime.Sub(hs.lastUpdate)
		if diff.Hours() >= 1 {
			return hs.Update()
		} else if currLen < hs.trendLength {
			return hs.Update()
		} else if len(*hs.trends) > 0 && ht.Count+h.Count > (*hs.trends)[len(*hs.trends)-1].Count {
			return hs.Update()
		}
	} else if currLen < hs.trendLength {
		return hs.Update()
	}
	return nil
}

func (hs *HashtagStore) GetHashtagByName(name string) (*model.Hashtag, error) {
	var h *model.Hashtag
	err := hs.db.FindOne(context.TODO(), bson.M{"name": name}).Decode(&h)
	return h, err
}

func (hs *HashtagStore) RemoveHashtag(name string) error {
	_, err := hs.db.DeleteOne(context.TODO(), bson.M{"name": name})
	if err != nil {
		return err
	}
	return nil
}

func (hs *HashtagStore) DeleteTweetHashtags(t *model.Tweet, hashtags map[string]int) error {
	for name, cnt := range hashtags {
		h, err := hs.GetHashtagByName(name)
		if err != nil {
			return err
		}
		newTweets := &[]primitive.ObjectID{}
		for _, id := range *h.Tweets {
			if id != t.ID {
				*newTweets = append(*newTweets, id)
			}
		}
		h.Count -= cnt
		h.Tweets = newTweets
		if h.Count == 0 {
			err := hs.RemoveHashtag(name)
			return err
		}
		_, err = hs.db.UpdateOne(context.TODO(), bson.M{"name": name}, bson.M{"$set": bson.M{"count": h.Count, "tweets": h.Tweets}})
		return err
	}
	return nil
}

func (hs *HashtagStore) GetHashtagTweets(name string) (*[]primitive.ObjectID, error) {
	h, err := hs.GetHashtagByName(name)
	if err != nil {
		return nil, err
	}
	return h.Tweets, nil
}

func (hs *HashtagStore) Update() error {
	opt := options.Find()
	opt.SetSort(bson.D{{"count", -1}})
	cur, err := hs.db.Find(context.TODO(), bson.D{}, opt)
	if err != nil {
		return err
	}
	hs.trends = &[]*model.Hashtag{}
	if err = cur.All(context.TODO(), hs.trends); err != nil {
		return err
	}
	fmt.Println(hs.trends)
	if len(*hs.trends) > 10 {
		*hs.trends = (*hs.trends)[:hs.trendLength]
	}
	hs.lastUpdate = time.Now()
	return nil
}

func (hs *HashtagStore) GetTrends() *[]*model.Hashtag {
	return hs.trends
}
