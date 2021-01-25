package store

import (
	"context"
	"fmt"
	"github.com/arman-aminian/twitter-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"regexp"
	"time"
)

type TweetStore struct {
	db *mongo.Collection
}

func NewTweetStore(db *mongo.Collection) *TweetStore {
	return &TweetStore{
		db: db,
	}
}

func (ts *TweetStore) CreateTweet(t *model.Tweet) error {
	_, err := ts.db.InsertOne(context.TODO(), t)
	return err
}

func (ts *TweetStore) RemoveTweet(t *model.Tweet) error {
	_, err := ts.db.DeleteOne(context.TODO(), t)
	return err
}

func (ts *TweetStore) GetTweetById(id *string) (*model.Tweet, error) {
	var t model.Tweet
	oid, err := primitive.ObjectIDFromHex(*id)
	if err != nil {
		return &t, nil
	}
	err = ts.db.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&t)
	return &t, err
}

func (ts *TweetStore) GetAllTweets() ([]bson.M, error) {
	var ret []bson.M
	cur, err := ts.db.Find(context.TODO(), bson.M{})
	if err = cur.All(context.TODO(), &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (ts *TweetStore) LikeTweet(t *model.Tweet, u *model.User) error {
	*t.Likes = append(*t.Likes, *model.NewOwner(u.Username, u.ProfilePicture))
	_, err := ts.db.UpdateOne(context.TODO(), bson.M{"_id": t.ID}, bson.M{"$set": bson.M{"likes": t.Likes}})
	if err != nil {
		return err
	}
	return nil
}

func (ts *TweetStore) UnLikeTweet(t *model.Tweet, u *model.User) error {
	newLikes := &[]model.Owner{}
	for _, o := range *t.Likes {
		if o.Username != u.Username {
			*newLikes = append(*newLikes, o)
		}
	}
	_, err := ts.db.UpdateOne(context.TODO(), bson.M{"_id": t.ID}, bson.M{"$set": bson.M{"likes": newLikes}})
	if err != nil {
		return err
	}
	t.Likes = newLikes
	return nil
}

func (ts *TweetStore) Retweet(t *model.Tweet, u *model.User) error {
	*t.Retweets = append(*t.Retweets, *model.NewOwner(u.Username, u.ProfilePicture))
	_, err := ts.db.UpdateOne(context.TODO(), bson.M{"_id": t.ID}, bson.M{"$set": bson.M{"retweets": t.Retweets}})
	if err != nil {
		return err
	}
	return nil
}

func (ts *TweetStore) UnRetweet(t *model.Tweet, u *model.User) error {
	newRetweets := &[]model.Owner{}
	for _, o := range *t.Retweets {
		if o.Username != u.Username {
			*newRetweets = append(*newRetweets, o)
		}
	}
	_, err := ts.db.UpdateOne(context.TODO(), bson.M{"_id": t.ID}, bson.M{"$set": bson.M{"retweets": newRetweets}})
	if err != nil {
		return err
	}
	t.Retweets = newRetweets
	return nil
}

func (ts *TweetStore) ExtractHashtags(t *model.Tweet) map[string]int {
	matchTags := regexp.MustCompile(`\B[#]\w*[a-zA-Z]+\w*`)
	res := map[string]int{}
	for _, v := range matchTags.FindAllString(t.Text, -1) {
		vn := v[1:]
		if _, ok := res[vn]; ok {
			res[vn] += 1
		} else {
			res[vn] = 1
		}
	}
	return res
}

func (ts *TweetStore) GetTimelineFromFollowingsUsernames(usernames []string) (*[]model.Tweet, error) {
	date := time.Now().Format("2006-01-02")
	var tweets []model.Tweet
	filter := bson.M{
		"$and": []bson.M{
			{"owner.username": bson.M{"$in": usernames}},
			{"date": date},
		},
	}
	// query := bson.M{"username": bson.M{"$in": usernames}, "date": bson.M{"$in": date}}
	res, err := ts.db.Find(context.TODO(), filter)
	if err != nil {
		println("injaaaaaa")
		fmt.Println(err)
		return nil, err
	}

	if err = res.All(context.TODO(), &tweets); err != nil {
		return nil, err
	}
	return &tweets, err
}
