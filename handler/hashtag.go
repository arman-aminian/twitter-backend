package handler

import (
	"errors"
	"github.com/arman-aminian/twitter-backend/model"
	"github.com/arman-aminian/twitter-backend/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (h *Handler) AddHashtag(name string, t *model.Tweet, count int) error {
	return h.hashtagStore.AddHashtag(&model.Hashtag{
		Name:   name,
		Tweets: &[]primitive.ObjectID{t.ID},
		Count:  count,
	})
}

type singleTrend struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type trendResponse struct {
	Trends    []singleTrend `json:"trends"`
	Timestamp time.Time     `json:"timestamp"`
}

func (h *Handler) GetTrends(c echo.Context) error {
	trends := h.hashtagStore.GetTrends()
	res := trendResponse{}
	for _, t := range *trends {
		s := singleTrend{
			Name:  t.Name,
			Count: t.Count,
		}
		res.Trends = append(res.Trends, s)
	}
	res.Timestamp = time.Now()
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) SearchHashtag(c echo.Context) error {
	query := &model.SearchQuery{}
	err := c.Bind(query)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if query.Query == "" {
		return c.JSON(http.StatusBadRequest, utils.NewError(errors.New("nothing to search for")))
	}
	result, err := h.hashtagStore.GetHashtagTweets(query.Query)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	tweets := &[]model.Tweet{}
	for _, id := range *result {
		s := id.Hex()
		t, err := h.tweetStore.GetTweetById(&s)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err))
		}
		*tweets = append(*tweets, *t)
	}
	return c.JSON(http.StatusOK, newTweetListResponse(stringFieldFromToken(c, "username"), tweets, len(*tweets)))
}
