package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
)

var clientInstance *mongo.Client
var clientInstanceError error
var mongoOnce sync.Once

const (
	// connect to local database
	// PATH = "mongodb://localhost:27017"

	// ### is your mongodb Atlas username
	// *** is your mongodb Atlas password
	// connect to mongodb Atlas
	PATH = "mongodb+srv://###:***@twitter-db.ugqqk.mongodb.net/twitter?retryWrites=true&w=majority"
)

func GetMongoClient() (*mongo.Client, error) {
	mongoOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(PATH)
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
		} else {
			err = client.Ping(context.TODO(), nil)
			if err != nil {
				clientInstanceError = err
			}
		}
		clientInstance = client
	})
	return clientInstance, clientInstanceError
}

func SetupUsersDb(mongoClient *mongo.Client) *mongo.Collection {
	usersDb := mongoClient.Database("twitter_db").Collection("users")
	createUniqueIndices(usersDb, "email")
	return usersDb
}

func SetupHashtagsDb(mongoClient *mongo.Client) *mongo.Collection {
	hashtagsDb := mongoClient.Database("twitter_db").Collection("hashtags")
	createUniqueIndices(hashtagsDb, "name")
	return hashtagsDb
}

func createUniqueIndices(db *mongo.Collection, field string) {
	_, err := db.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: field, Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func SetupTweetsDb(mongoClient *mongo.Client) *mongo.Collection {
	tweetsDb := mongoClient.Database("twitter_db").Collection("tweets")
	_, err := tweetsDb.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "text", Value: "text"}},
			Options: nil,
		})

	if err != nil {
		log.Fatal(err)
	}
	return tweetsDb
}
