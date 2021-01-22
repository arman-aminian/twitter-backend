package model

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            primitive.ObjectID    `bson:"_id,omitempty"`
	Username      string                `json:"username" bson:"username"`
	Email         string                `json:"email" bson:"email"`
	Password      string                `json:"password" bson:"password"`
	Tweets        *[]primitive.ObjectID `json:"tweets" bson:"tweets"`
	Followings    *[]primitive.ObjectID `json:"followings" bson:"followings"`
	Followers     *[]primitive.ObjectID `json:"followers" bson:"followers"`
	Notifications *[]Event              `json:"notifications" bson:"notifications"`
	Logs          *[]Event              `json:"logs" bson:"logs"`
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
