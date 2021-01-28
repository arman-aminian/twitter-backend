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
		Username       string `json:"username" bson:"_id"`
		Email          string `json:"email"`
		Name           string `json:"name"`
		Bio            string `json:"bio"`
		ProfilePicture string `json:"profile_picture"`
		Token          string `json:"token"`
	} `json:"user"`
}

func newUserResponse(u *model.User) *userResponse {
	r := new(userResponse)
	r.User.Username = u.Username
	r.User.Email = u.Email
	r.User.Name = u.Name
	r.User.Bio = u.Bio
	r.User.ProfilePicture = u.ProfilePicture
	r.User.Token = utils.GenerateJWT(u.Username)
	return r
}

type profileResponse struct {
	Profile struct {
		Name           string                `json:"name"`
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
	r.Profile.Name = u.Name
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
	temp := *u.Followers
	for i := range temp {
		temp[i].IsFollowing, _ = us.IsFollower(temp[i].Username, srcUsername)
	}
	l.Followers = &temp
	temp2 := *u.Followings
	for i := range temp2 {
		temp2[i].IsFollowing, _ = us.IsFollower(temp2[i].Username, srcUsername)
	}
	l.Followings = &temp2
	l.IsFollowing, _ = us.IsFollower(u.Username, srcUsername)
	return l
}

type OwnerListResponse struct {
	Users *[]model.Owner `json:"users"`
}

func newOwnerList(us user.Store, srcUsername string, users *[]model.Owner) *OwnerListResponse {
	l := new(OwnerListResponse)

	temp := *users
	for i := range temp {
		temp[i].IsFollowing, _ = us.IsFollower(temp[i].Username, srcUsername)
	}
	l.Users = &temp
	return l
}

type SingleEventResponse struct {
	Mode      string      `json:"mode"`
	Source    model.Owner `json:"source"`
	Target    model.Owner `json:"target"`
	Content   string      `json:"content"`
	Tweet     model.Tweet `json:"tweet"`
	TimeStamp time.Time   `json:"timestamp"`
}

type EventListResponse struct {
	Events []SingleEventResponse `json:"events"`
}

func newLogsList(u *model.User) *EventListResponse {
	ret := new(EventListResponse)
	for _, e := range *u.Logs {
		l := new(SingleEventResponse)
		l.Mode = e.Mode
		l.Source = e.Source
		l.Target = e.Target
		l.Content = e.Content
		l.Tweet = *e.Tweet
		l.TimeStamp = e.TimeStamp
		ret.Events = append([]SingleEventResponse{*l}, ret.Events...)
	}
	return ret
}

func newNotificationsList(u *model.User) *EventListResponse {
	ret := new(EventListResponse)
	for _, e := range *u.Notifications {
		l := new(SingleEventResponse)
		l.Mode = e.Mode
		l.Source = e.Source
		l.Target = e.Target
		l.Content = e.Content
		l.Tweet = *e.Tweet
		l.TimeStamp = e.TimeStamp
		ret.Events = append([]SingleEventResponse{*l}, ret.Events...)
	}
	return ret
}

//	********************** Tweet Response **********************

type tweetResponse struct {
	ID            string                `json:"id"`
	Text          string                `json:"text"`
	Media         string                `json:"media"`
	Liked         bool                  `json:"liked"`
	LikesCount    int                   `json:"likes_count"`
	Retweeted     bool                  `json:"retweeted"`
	RetweetsCount int                   `json:"retweets_count"`
	Time          time.Time             `json:"time" bson:"time"`
	Owner         model.Owner           `json:"owner"`
	Parents       *[]model.CommentTweet `json:"parents" bson:"parents"`
	Comments      *[]model.CommentTweet `json:"comments" bson:"comments"`
}

type singleTweetResponse struct {
	Tweet *tweetResponse `json:"tweet"`
}

type tweetsResponse struct {
	Tweets []tweetResponse `json:"tweets"`
}

type tweetListResponse struct {
	Tweets      []tweetResponse `json:"tweets"`
	TweetsCount int             `json:"tweetsCount"`
}

func newTweetsResponse(username string, tweets *[]model.Tweet) *tweetsResponse {
	tr := make([]tweetResponse, len(*tweets))
	if tweets == nil {
		return &tweetsResponse{tr}
	}
	for i, tweet := range *tweets {
		tr[i].ID = tweet.ID.Hex()
		tr[i].Text = tweet.Text
		tr[i].Parents = tweet.Parents
		tr[i].Comments = tweet.Comments
		tr[i].Media = tweet.Media
		tr[i].Time = tweet.Time
		
		for _, t := range *tweet.Likes {
			if t.Username == username {
				tr[i].Liked = true
				break
			}
		}
		tr[i].LikesCount = len(*tweet.Likes)
		for _, t := range *tweet.Retweets {
			if t.Username == username {
				tr[i].Retweeted = true
				break
			}
		}
		tr[i].RetweetsCount = len(*tweet.Retweets)
		tr[i].Owner.Username = tweet.Owner.Username
		tr[i].Owner.ProfilePicture = tweet.Owner.ProfilePicture
	}
	return &tweetsResponse{tr}
}

func newTweetResponse(c echo.Context, t *model.Tweet) *singleTweetResponse {
	tr := new(tweetResponse)
	tr.ID = t.ID.Hex()
	tr.Text = t.Text
	tr.Parents = t.Parents
	tr.Comments = t.Comments
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

func newTweetListResponse(c echo.Context, username string, tweets *[]model.Tweet, size int) *tweetListResponse {
	tr := make([]tweetResponse, size)
	if tweets == nil {
		return &tweetListResponse{tr, size}
	}
	for i, tweet := range *tweets {
		tr[i].ID = tweet.ID.Hex()
		tr[i].Text = tweet.Text
		tr[i].Parents = tweet.Parents
		tr[i].Comments = tweet.Comments
		tr[i].Media = tweet.Media
		tr[i].Time = tweet.Time
		
		for _, t := range *tweet.Likes {
			if t.Username == username {
				tr[i].Liked = true
				break
			}
		}
		tr[i].LikesCount = len(*tweet.Likes)
		for _, t := range *tweet.Retweets {
			if t.Username == username {
				tr[i].Retweeted = true
				break
			}
		}
		tr[i].RetweetsCount = len(*tweet.Retweets)
		tr[i].Owner.Username = tweet.Owner.Username
		tr[i].Owner.ProfilePicture = tweet.Owner.ProfilePicture
	}
	return &tweetListResponse{tr, size}
}

type tweetLikeAndRetweetResponse struct {
	LikesList    *[]model.Owner `json:"likes" bson:"likes"`
	RetweetsList *[]model.Owner `json:"retweets" bson:"retweets"`
}

func newLikeAndRetweetResponse(us user.Store, srcUsername string, t *model.Tweet) *tweetLikeAndRetweetResponse {
	tr := new(tweetLikeAndRetweetResponse)

	temp := *t.Likes
	for i := range temp {
		temp[i].IsFollowing, _ = us.IsFollower(temp[i].Username, srcUsername)
	}
	tr.LikesList = &temp

	temp2 := *t.Retweets
	for i := range temp2 {
		temp2[i].IsFollowing, _ = us.IsFollower(temp2[i].Username, srcUsername)
	}
	tr.RetweetsList = &temp2
	return tr
}

type timelineResponse struct {
	timeline *[]model.Tweet
}
