package handler

import (
	"github.com/arman-aminian/twitter-backend/model"
	"github.com/arman-aminian/twitter-backend/user"
	"github.com/arman-aminian/twitter-backend/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
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
		Followings     *[]model.Owner        `json:"followings"`
		Followers      *[]model.Owner        `json:"followers"`
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
	r.Profile.IsFollowing, _ = us.IsFollower(u.Username, srcUsername)
	return r
}

type FollowersAndFollowingListResponse struct {
	Followers   *[]model.Owner `json:"followers" bson:"followers"`
	Followings  *[]model.Owner `json:"followings" bson:"followings"`
	IsFollowing bool           `json:"is_following, omitempty"`
}

func newFollowingAndFollowersList(us user.Store, srcUsername string, u *model.User) *FollowersAndFollowingListResponse {
	l := new(FollowersAndFollowingListResponse)
	l.Followers = u.Followers
	l.Followings = u.Followings
	l.IsFollowing, _ = us.IsFollower(u.Username, srcUsername)
	return l
}

//	********************** Tweet Response **********************

type tweetResponse struct {
	ID            string      `json:"id"`
	Text          string      `json:"text"`
	Media         string      `json:"media"`
	Liked         bool        `json:"liked"`
	LikesCount    int         `json:"likes_count"`
	Retweeted     bool        `json:"retweeted"`
	RetweetsCount int         `json:"retweets_count"`
	Time          time.Time   `json:"time" bson:"time"`
	Owner         model.Owner `json:"owner"`
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
	tr.ID = t.ID.Hex()
	tr.Text = t.Text
	tr.Media = t.Media
	tr.Time = t.Time
	for _, u := range *t.Likes {
		if u.Username == stringFieldFromToken(c, "username") {
			tr.Liked = true
			break
		}
	}
	tr.LikesCount = len(*t.Likes)
	for _, u := range *t.Retweets {
		if u.Username == stringFieldFromToken(c, "username") {
			tr.Retweeted = true
			break
		}
	}
	tr.RetweetsCount = len(*t.Retweets)
	tr.Owner.Username = t.Owner.Username
	tr.Owner.ProfilePicture = t.Owner.ProfilePicture

	return &singleTweetResponse{tr}
}

type tweetLikeAndRetweetResponse struct {
	LikesList    *[]model.Owner `json:"likes" bson:"likes"`
	RetweetsList *[]model.Owner `json:"retweets" bson:"retweets"`
}

func newLikeAndRetweetResponse(t *model.Tweet) *tweetLikeAndRetweetResponse {
	tr := new(tweetLikeAndRetweetResponse)
	tr.LikesList = t.Likes
	tr.RetweetsList = t.Retweets
	return tr
}

type timelineResponse struct {
	timeline *[]model.Tweet
}
