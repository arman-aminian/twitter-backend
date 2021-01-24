package handler

import (
	"github.com/arman-aminian/twitter-backend/model"
	"github.com/arman-aminian/twitter-backend/user"
	"github.com/arman-aminian/twitter-backend/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userResponse struct {
	User struct {
		Username string `json:"username" bson:"_id"`
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
		IsFollowing    bool                  `json:"is_following, omitempty"`
		Username       string                `json:"username" bson:"_id"`
		Bio            string                `json:"bio"`
		ProfilePicture string                `json:"profile_picture"`
		HeaderPicture  string                `json:"header_picture"`
		Tweets         *[]primitive.ObjectID `json:"tweets"`
		Followings     *[]string             `json:"followings"`
		Followers      *[]string             `json:"followers"`
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

//	********************** Tweet Response **********************

type tweetResponse struct {
	Text          string `json:"text"`
	Media         string `json:"media"`
	Liked         bool   `json:"liked"`
	LikesCount    int    `json:"likes_count"`
	Retweeted     bool   `json:"retweeted"`
	RetweetsCount int    `json:"retweets_count"`
	Owner         struct {
		Username       string `json:"username"`
		ProfilePicture string `json:"profile_picture"`
	} `json:"owner"`
}

type singleTweetResponse struct {
	Tweet *tweetResponse `json:"tweet"`
}

type tweetListResponse struct {
	Tweets      []*tweetResponse `json:"tweets"`
	TweetsCount int              `json:"tweetsCount"`
}

func newTweetResponse(c echo.Context, t *model.Tweet) *singleTweetResponse {
	tr := new(tweetResponse)
	tr.Text = t.Text
	tr.Media = t.Media
	for _, u := range *t.Likes {
		if u == usernameFromToken(c) {
			tr.Liked = true
		}
	}
	tr.LikesCount = len(*t.Likes)
	for _, u := range *t.Retweets {
		if u == usernameFromToken(c) {
			tr.Retweeted = true
		}
	}
	tr.RetweetsCount = len(*t.Retweets)
	tr.Owner.Username = t.Owner.Username
	tr.Owner.ProfilePicture = t.Owner.ProfilePicture

	return &singleTweetResponse{tr}
}
