package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Username      string                `json:"username" bson:"username"`
	Email         string                `json:"email" bson:"email"`
	Password      string                `json:"password" bson:"password"`
	Tweets        *[]primitive.ObjectID `json:"tweets" bson:"tweets"`
	Followings    *[]primitive.ObjectID `json:"followings" bson:"followings"`
	Followers     *[]primitive.ObjectID `json:"followers" bson:"followers"`
	Notifications *[]Event              `json:"notifications" bson:"notifications"`
	Logs          *[]Event              `json:"logs" bson:"logs"`
}
