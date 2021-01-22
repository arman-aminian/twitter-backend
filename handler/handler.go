package handler

import "github.com/arman-aminian/twitter-backend/user"

type Handler struct {
	userStore user.Store
}

func NewHandler(us user.Store) *Handler {
	return &Handler{
		userStore: us,
	}
}
