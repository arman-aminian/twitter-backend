package store

import (
	"context"
	"github.com/arman-aminian/twitter-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore struct {
	db *mongo.Collection
}

func NewUserStore(db *mongo.Collection) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (us *UserStore) Create(u *model.User) error {
	_, err := us.db.InsertOne(context.TODO(), u)
	return err
}

func (us *UserStore) Remove(username string) error {
	_, err := us.db.DeleteOne(context.TODO(), bson.M{"_id": username})
	return err
}

func (us *UserStore) Update(u *model.User) error {
	_, err := us.db.UpdateOne(context.TODO(),
		bson.M{"_id": u.User},
		bson.M{"$set": bson.M{
			"username": u.Username,
			"email":    u.Email,
			"password": u.Password,
		},
		})
	return err
}

func (us *UserStore) UpdateProfile(u *model.User) error {
	_, err := us.db.UpdateOne(context.TODO(),
		bson.M{"username": u.Username},
		bson.M{"$set": bson.M{
			"bio":             u.Bio,
			"profile_picture": u.ProfilePicture,
			"header_picture":  u.HeaderPicture,
		},
		})
	return err
}

func (us *UserStore) GetByEmail(email string) (*model.User, error) {
	var u model.User
	err := us.db.FindOne(context.TODO(), bson.M{"email": email}).Decode(&u)
	return &u, err
}

func (us *UserStore) GetByUsername(username string) (*model.User, error) {
	var u model.User
	err := us.db.FindOne(context.TODO(), bson.M{"_id": username}).Decode(&u)
	return &u, err
}

func (us *UserStore) AddFollower(u *model.User, follower *model.User) error {
	*u.Followers = append(*u.Followers, follower.Username)
	_, err := us.db.UpdateOne(context.TODO(), bson.M{"_id": u.Username}, bson.M{"$set": bson.M{"followers": u.Followers}})
	if err != nil {
		return err
	}
	*follower.Followings = append(*follower.Followings, u.Username)
	_, err = us.db.UpdateOne(context.TODO(), bson.M{"_id": follower.Username}, bson.M{"$set": bson.M{"followings": follower.Followings}})
	if err != nil {
		return err
	}
	return nil
}

func (us *UserStore) IsFollower(username, followerUsername string) (bool, error) {
	u, err := us.GetByUsername(username)
	if err != nil {
		return false, err
	}
	follower, err := us.GetByUsername(followerUsername)
	if err != nil {
		return false, nil
	}
	doesFollow := false
	for _, f := range *u.Followers {
		if f == follower.Username {
			doesFollow = true
			break
		}
	}
	hasInFollowings := false
	for _, f := range *follower.Followings {
		if f == u.Username {
			hasInFollowings = true
			break
		}
	}
	return doesFollow && hasInFollowings, nil
}

func (us *UserStore) AddTweet(u *model.User, t *model.Tweet) error {
	*u.Tweets = append(*u.Tweets, t.ID)
	_, err := us.db.UpdateOne(context.TODO(), bson.M{"_id": u.Username}, bson.M{"$set": bson.M{"tweets": u.Tweets}})
	return err
}
