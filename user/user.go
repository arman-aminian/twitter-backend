package user

import (
	"github.com/arman-aminian/twitter-backend/model"
)

type Store interface {
	Create(*model.User) error
	Remove(field, value string) error
	Update(old *model.User, new *model.User) error
	UpdateProfile(u *model.User) error

	GetByEmail(string) (*model.User, error)
	GetByUsername(string) (*model.User, error)
	AddFollower(u *model.User, follower *model.User) error
	RemoveFollower(u *model.User, follower *model.User) error
	IsFollower(username, followerUsername string) (bool, error)

	AddTweet(u *model.User, t *model.Tweet) error

	GetUserListFromUsernameList(usernames []string) (*[]model.User, error)
}
