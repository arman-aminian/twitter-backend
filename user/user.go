package user

import "github.com/arman-aminian/twitter-backend/model"

type Store interface {
	Create(*model.User) error

	GetByEmail(string) (*model.User, error)
	GetByUsername(string) (*model.User, error)
}
