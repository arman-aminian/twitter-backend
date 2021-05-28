package store

import (
	"context"
	"github.com/arman-aminian/twitter-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (us *UserStore) Remove(field, value string) error {
	_, err := us.db.DeleteOne(context.TODO(), bson.M{field: value})
	return err
}

func (us *UserStore) Update(old *model.User, new *model.User) error {
	var err error
	if old.Username != new.Username {
		err = us.Remove("_id", old.Username)
	} else if old.Email != new.Email {
		err = us.Remove("email", old.Email)
	} else if old.Password != new.Password {
		err = us.Remove("password", old.Password)
	}
	err = us.Create(new)
	return err
}

func (us *UserStore) UpdateProfile(u *model.User) error {
	_, err := us.db.UpdateOne(context.TODO(),
		bson.M{"_id": u.Username},
		bson.M{"$set": bson.M{
			"name":            u.Name,
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
	*u.Followers = append(*u.Followers, *model.NewOwner(follower.Username, follower.ProfilePicture, follower.Name, follower.Bio))
	_, err := us.db.UpdateOne(context.TODO(), bson.M{"_id": u.Username}, bson.M{"$set": bson.M{"followers": u.Followers}})
	if err != nil {
		return err
	}
	*follower.Followings = append(*follower.Followings, *model.NewOwner(u.Username, u.ProfilePicture, u.Name, u.Bio))
	_, err = us.db.UpdateOne(context.TODO(), bson.M{"_id": follower.Username}, bson.M{"$set": bson.M{"followings": follower.Followings}})
	if err != nil {
		return err
	}
	return nil
}

func (us *UserStore) RemoveFollower(u *model.User, follower *model.User) error {
	newFollowers := &[]model.Owner{}
	for _, o := range *u.Followers {
		if o.Username != follower.Username {
			*newFollowers = append(*newFollowers, o)
		}
	}
	_, err := us.db.UpdateOne(context.TODO(), bson.M{"_id": u.Username}, bson.M{"$set": bson.M{"followers": newFollowers}})
	if err != nil {
		return err
	}
	u.Followers = newFollowers

	newFollowings := &[]model.Owner{}
	for _, o := range *follower.Followings {
		if o.Username != u.Username {
			*newFollowings = append(*newFollowings, o)
		}
	}
	_, err = us.db.UpdateOne(context.TODO(), bson.M{"_id": follower.Username}, bson.M{"$set": bson.M{"followings": newFollowings}})
	if err != nil {
		return err
	}
	follower.Followings = newFollowings
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
	for _, o := range *u.Followers {
		if o.Username == follower.Username {
			doesFollow = true
			break
		}
	}
	hasInFollowings := false
	for _, o := range *follower.Followings {
		if o.Username == u.Username {
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

func (us *UserStore) RemoveTweet(u *model.User, id *string) error {
	oid, err := primitive.ObjectIDFromHex(*id)
	if err != nil {
		return err
	}
	newTweets := &[]primitive.ObjectID{}
	for _, tid := range *u.Tweets {
		if tid != oid {
			*newTweets = append(*newTweets, tid)
		}
	}
	_, err = us.db.UpdateOne(context.TODO(), bson.M{"_id": u.Username}, bson.M{"$set": bson.M{"tweets": newTweets}})
	return err
}

func (us *UserStore) AddLog(u *model.User, e *model.Event) error {
	*u.Logs = append(*u.Logs, *e)
	_, err := us.db.UpdateOne(context.TODO(), bson.M{"_id": u.Username}, bson.M{"$set": bson.M{"logs": u.Logs}})
	if err != nil {
		return err
	}
	return nil
}

func (us *UserStore) AddNotification(u *model.User, e *model.Event) error {
	*u.Notifications = append(*u.Notifications, *e)
	_, err := us.db.UpdateOne(context.TODO(), bson.M{"_id": u.Username}, bson.M{"$set": bson.M{"notifications": u.Notifications}})
	if err != nil {
		return err
	}
	return nil
}

func (us *UserStore) GetUserListFromUsernameList(usernames []string) (*[]model.User, error) {
	var users []model.User
	query := bson.M{"_id": bson.M{"$in": usernames}}
	res, err := us.db.Find(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	if err = res.All(context.TODO(), &users); err != nil {
		return nil, err
	}
	return &users, err
}

func (us *UserStore) GetTweetIdListFromUsernameList(usernames []string) (*[]primitive.ObjectID, error) {
	var users []model.User
	query := bson.M{"_id": bson.M{"$in": usernames}}
	res, err := us.db.Find(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	if err = res.All(context.TODO(), &users); err != nil {
		return nil, err
	}
	var tweetsId []primitive.ObjectID
	for _, user := range users {
		tweetsId = append(tweetsId, *user.Tweets...)
	}
	return &tweetsId, err
}

func (us *UserStore) GetUsernameSearchResult(username string) (*[]model.Owner, error) {
	var users []model.User
	reg := "^" + username // usernames that starts with "query"
	query := bson.M{"_id": bson.M{"$regex": reg}}
	res, err := us.db.Find(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	if err = res.All(context.TODO(), &users); err != nil {
		return nil, err
	}
	var result []model.Owner
	for _, user := range users {
		result = append(result, model.Owner{
			Username:       user.Username,
			ProfilePicture: user.ProfilePicture,
			Name:           user.Name,
			Bio:            user.Bio,
		})
	}
	return &result, err
}
