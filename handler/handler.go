package handler

import (
	"github.com/arman-aminian/twitter-backend/tweet"
	"github.com/arman-aminian/twitter-backend/user"
)

type Handler struct {
	userStore  user.Store
	tweetStore tweet.Store
}

func NewHandler(us user.Store, ts tweet.Store) *Handler {
	return &Handler{
		userStore:  us,
		tweetStore: ts,
	}
}
