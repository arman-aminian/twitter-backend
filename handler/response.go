package handler

import (
	"github.com/arman-aminian/twitter-backend/model"
	"github.com/arman-aminian/twitter-backend/user"
	"github.com/arman-aminian/twitter-backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userResponse struct {
	User struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Token    string `json:"token"`
	} `json:"user"`
}

func newUserResponse(u *model.User) *userResponse {
	r := new(userResponse)
	r.User.Username = u.Username
	r.User.Email = u.Email
	r.User.Token = utils.GenerateJWT(u.Username)
	return r
}

type profileResponse struct {
	Profile struct {
		IsFollowing    bool                  `json:"is_following"`
		Username       string                `json:"username"`
		Bio            string                `json:"bio"`
		ProfilePicture string                `json:"profile_picture"`
		HeaderPicture  string                `json:"header_picture"`
		Tweets         *[]primitive.ObjectID `json:"tweets"`
		Followings     *[]primitive.ObjectID `json:"followings"`
		Followers      *[]primitive.ObjectID `json:"followers"`
	} `json:"profile"`
}

func newProfileResponse(us user.Store, srcUsername string, u *model.User) *profileResponse {
	r := new(profileResponse)
	r.Profile.Username = u.Username
	r.Profile.Bio = u.Bio
	r.Profile.ProfilePicture = u.ProfilePicture
	r.Profile.HeaderPicture = u.HeaderPicture
	r.Profile.Tweets = u.Tweets
	r.Profile.Followings = u.Followings
	r.Profile.Followers = u.Followers
	// Does srcUsername follow u.Username?
	r.Profile.IsFollowing, _ = us.IsFollower(u.Username, srcUsername)
	return r
}
