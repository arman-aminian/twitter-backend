package tweet

import "github.com/arman-aminian/twitter-backend/model"

type Store interface {
	CreateTweet(*model.Tweet) error
}
