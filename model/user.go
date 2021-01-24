package model

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username" bson:"_id"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`

	// Fluff shown to user as profile
	Bio            string `json:"bio" bson:"bio"`
	ProfilePicture string `json:"profile_picture" bson:"profile_picture"`
	HeaderPicture  string `json:"header_picture" bson:"header_picture"`

	Tweets        *[]primitive.ObjectID `json:"tweets" bson:"tweets"`
	Followings    *[]Owner              `json:"followings" bson:"followings"`
	Followers     *[]Owner              `json:"followers" bson:"followers"`
	Notifications *[]Event              `json:"notifications" bson:"notifications"`
	Logs          *[]Event              `json:"logs" bson:"logs"`
}

func NewUser() *User {
	var u User
	u.Tweets = &[]primitive.ObjectID{}
	u.Followings = &[]Owner{}
	u.Followers = &[]Owner{}
	u.Notifications = &[]Event{}
	u.Logs = &[]Event{}
	return &u
}

func (u *User) HashPassword(plain string) (string, error) {
	if len(plain) == 0 {
		return "", errors.New("password should not be empty")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(h), err
}

func (u *User) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	return err == nil
}
