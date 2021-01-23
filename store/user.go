package store

import (
	"context"
	"fmt"
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
	_, err := us.db.DeleteOne(context.TODO(), bson.M{"username": username})
	return err
}

func (us *UserStore) GetByEmail(email string) (*model.User, error) {
	var u model.User
	err := us.db.FindOne(context.TODO(), bson.M{"email": email}).Decode(&u)
	return &u, err
}

func (us *UserStore) GetByUsername(username string) (*model.User, error) {
	var u model.User
	err := us.db.FindOne(context.TODO(), bson.M{"username": username}).Decode(&u)
	return &u, err
}

func (us *UserStore) AddFollower(u *model.User, follower *model.User) error {
	fmt.Printf("before: %v %v\n", *u.Followers, follower)
	*u.Followers = append(*u.Followers, follower.ID)
	fmt.Printf("after: %v\n", *u.Followers)
	*follower.Followings = append(*follower.Followings, u.ID)
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
		if f == follower.ID {
			doesFollow = true
			break
		}
	}
	hasInFollowings := false
	for _, f := range *follower.Followings {
		if f == u.ID {
			hasInFollowings = true
			break
		}
	}
	return doesFollow && hasInFollowings, nil
}
