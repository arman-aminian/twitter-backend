package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hashtag struct {
	Name   string                `json:"name" bson:"name"`
	Tweets *[]primitive.ObjectID `json:"tweets" bson:"tweets"`
}
