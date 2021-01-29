package handler

import (
	"errors"
	"github.com/arman-aminian/twitter-backend/model"
	"github.com/arman-aminian/twitter-backend/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"os"
	"sort"
	"time"
)

// CreateTweet godoc
// @Summary Create an tweet
// @Description Create an tweet
// @ID create-tweet
// @Tags tweet
// @Accept  json
// @Produce  json
// @Param tweet body tweetCreateRequest true "Tweet to create made of text and media"
// @Success 201 {object} singleTweetResponse
// @Failure 404 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /tweets [post]
func (h *Handler) CreateTweet(c echo.Context) error {
	t := model.NewTweet()

	t.ID = primitive.NewObjectID()
	t.Text = c.FormValue("text")

	file, err := c.FormFile("media")
	if err == nil {
		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err))
		}
		defer src.Close()

		mediaFolderName := "media/tweet-assets/"
		mediaPath := mediaFolderName + file.Filename
		dst, err := os.Create(mediaPath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err))
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err))
		}
		t.Media = mediaPath
	} else {
		// Tweet without media
		t.Media = ""
	}

	u, err := h.userStore.GetByUsername(stringFieldFromToken(c, "username"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	t.Owner.Username = u.Username
	t.Owner.ProfilePicture = u.ProfilePicture
	t.Time = time.Now()
	t.Date = time.Now().Format("2006-01-02")

	parentId := c.FormValue("parent")
	if parentId != "" {
		p, err := h.tweetStore.GetTweetById(&parentId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err))
		}
		if p == nil {
			return c.JSON(http.StatusNotFound, utils.NotFound())
		}

		pid, _ := h.tweetStore.GetTweetById(&parentId)
		//
		temp := *model.NewCommentTweet(*pid)
		temp.CommentsCount = temp.CommentsCount + 1
		*t.Parents = append(*p.Parents, temp)
		err = h.tweetStore.AddCommentToTweet(p, model.NewCommentTweet(*t))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err))
		}
	}

	err = h.tweetStore.CreateTweet(t)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	err = h.userStore.AddTweet(u, t)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	hashtags := h.tweetStore.ExtractHashtags(t)
	for name, cnt := range hashtags {
		err = h.AddHashtag(name, t, cnt)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
		}
	}

	res := newTweetResponse(c, t)
	return c.JSON(http.StatusCreated, res)
}

func (h *Handler) GetTweetAssetFile(c echo.Context) error {
	mediaFolderName := "media/tweet-assets/"
	mediaPath := mediaFolderName + c.Param("filename")
	return c.File(mediaPath)
}

// GetTweet godoc
// @Description Create an tweet. Auth is optional.
// @ID get-tweet
// @Tags tweet
// @Produce  json
// @Success 201 {object} singleTweetResponse
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /tweets/{id} [get]
func (h *Handler) GetTweet(c echo.Context) error {
	id := c.Param("id")
	t, err := h.tweetStore.GetTweetById(&id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if t == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	return c.JSON(http.StatusOK, newTweetResponse(c, t))
}

// GetTweets godoc
// @Description Get all of the tweets of a user. Auth is required.
// @ID get-tweet
// @Tags tweet
// @Produce  json
// @Success 201 {object} singleTweetResponse
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /tweets/{id} [post]
func (h *Handler) GetTweets(c echo.Context) error {
	tweets := &model.TweetIdList{}
	err := c.Bind(tweets)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if len(tweets.Tweets) == 0 {
		return c.JSON(http.StatusBadRequest, utils.NewError(errors.New("nothing to search for")))
	}
	res, err := h.tweetStore.GetTweets(tweets.Tweets)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	// sort tweets
	sorted := *res
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Time.After(sorted[j].Time)
	})
	return c.JSON(http.StatusOK, newTweetsResponse(stringFieldFromToken(c, "username"), &sorted))
}

// DeleteTweet godoc
// @Description Delete a tweet from a user's tweets based on the token. Auth is required.
// @ID delete-tweet
// @Tags tweet
// @Produce  json
// @Success 201 {object} singleTweetResponse
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /tweets/{id} [delete]
func (h *Handler) DeleteTweet(c echo.Context) error {
	id := c.Param("id")
	t, err := h.tweetStore.GetTweetById(&id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if t == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	u, err := h.userStore.GetByUsername(t.Owner.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	if u.Username != t.Owner.Username {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("you can only delete you tweets")))
	}

	err = h.userStore.RemoveTweet(u, &id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	err = h.tweetStore.RemoveTweet(t)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	hashtags := h.tweetStore.ExtractHashtags(t)
	_ = h.hashtagStore.DeleteTweetHashtags(t, hashtags)

	temp := *t.Parents
	parent := temp[len(*t.Parents)-1]

	pid := parent.ID.Hex()
	pt, err := h.tweetStore.GetTweetById(&pid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	err = h.tweetStore.RemoveComment(pt, &t.ID)
	if err != nil {
		panic(err)
	}
	return c.JSON(http.StatusOK, newTweetResponse(c, t))
}

// GetTweetLikeAndRetweetList godoc
// @Description Get the list of users who liked and retweeted this tweet. Auth not required.
// @ID get-likes-retweets-list
// @Tags tweet
// @Produce  json
// @Param id path string true "Id of the tweet to get the list from."
// @Success 200 {object} tweetLikeAndRetweetResponse
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /tweets/{id}/list [get]
func (h *Handler) GetTweetLikeAndRetweetList(c echo.Context) error {
	id := c.Param("id")
	t, err := h.tweetStore.GetTweetById(&id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if t == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	return c.JSON(http.StatusOK, newLikeAndRetweetResponse(h.userStore, stringFieldFromToken(c, "username"), t))
}

// Like godoc
// @Description Like a tweet. Auth is required.
// @ID like
// @Tags like
// @Produce  json
// @Param id path string true "id of the article that you want to like"
// @Success 200 {object} singleTweetResponse
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /tweets/{id}/like [post]
func (h *Handler) Like(c echo.Context) error {
	id := c.Param("id")
	t, err := h.tweetStore.GetTweetById(&id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if t == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	u, err := h.userStore.GetByUsername(stringFieldFromToken(c, "username"))
	if u == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	for _, o := range *t.Likes {
		if u.Username == o.Username {
			return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("already liked")))
		}
	}

	if err := h.tweetStore.LikeTweet(t, u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	e := h.CreateLikeEvent(u, t)
	err = h.userStore.AddLog(u, e)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	target, err := h.userStore.GetByUsername(t.Owner.Username)
	if target == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	err = h.userStore.AddNotification(target, e)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, newTweetResponse(c, t))
}

// UnLike godoc
// @Description UnLike a tweet. Auth is required.
// @ID unlike
// @Tags like
// @Produce  json
// @Param id path string true "id of the article that you want to unlike"
// @Success 200 {object} singleTweetResponse
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /tweets/{id}/like [delete]
func (h *Handler) UnLike(c echo.Context) error {
	id := c.Param("id")
	t, err := h.tweetStore.GetTweetById(&id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if t == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	u, err := h.userStore.GetByUsername(stringFieldFromToken(c, "username"))
	if u == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	b := false
	for _, o := range *t.Likes {
		if u.Username == o.Username {
			b = true
		}
	}
	if !b {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("hasn't liked")))
	}

	if err := h.tweetStore.UnLikeTweet(t, u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, newTweetResponse(c, t))
}

// Retweet godoc
// @Description retweet a tweet. Auth is required.
// @ID retweet
// @Tags retweet
// @Produce  json
// @Param id path string true "id of the article that you want to retweet"
// @Success 200 {object} singleTweetResponse
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /tweets/{id}/retweet [post]
func (h *Handler) Retweet(c echo.Context) error {
	id := c.Param("id")
	t, err := h.tweetStore.GetTweetById(&id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if t == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	u, err := h.userStore.GetByUsername(stringFieldFromToken(c, "username"))
	if u == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	for _, o := range *t.Retweets {
		if u.Username == o.Username {
			return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("already retweeted")))
		}
	}

	if err := h.tweetStore.Retweet(t, u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	err = h.userStore.AddTweet(u, t)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	e := h.CreateRetweetEvent(u, t)
	err = h.userStore.AddLog(u, e)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	target, err := h.userStore.GetByUsername(t.Owner.Username)
	if target == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	err = h.userStore.AddNotification(target, e)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, newTweetResponse(c, t))
}

// UnRetweet godoc
// @Description UnRetweet a tweet. Auth is required.
// @ID unretweet
// @Tags unretweet
// @Produce  json
// @Param id path string true "id of the article that you want to unretweet"
// @Success 200 {object} singleTweetResponse
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /tweets/{id}/retweet [delete]
func (h *Handler) UnRetweet(c echo.Context) error {
	id := c.Param("id")
	t, err := h.tweetStore.GetTweetById(&id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if t == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	u, err := h.userStore.GetByUsername(stringFieldFromToken(c, "username"))
	if u == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	b := false
	for _, o := range *t.Retweets {
		if u.Username == o.Username {
			b = true
		}
	}
	if !b {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("hasn't retweeted")))
	}

	if err := h.tweetStore.UnRetweet(t, u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	err = h.userStore.RemoveTweet(u, &id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, newTweetResponse(c, t))
}
