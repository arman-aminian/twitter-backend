package handler

import (
	"github.com/arman-aminian/twitter-backend/hashtag"
	"github.com/arman-aminian/twitter-backend/tweet"
	"github.com/arman-aminian/twitter-backend/user"
)

type Handler struct {
	userStore    user.Store
	tweetStore   tweet.Store
	hashtagStore hashtag.Store
}

func NewHandler(us user.Store, ts tweet.Store, hh hashtag.Store) *Handler {
	return &Handler{
		userStore:    us,
		tweetStore:   ts,
		hashtagStore: hh,
	}
}
