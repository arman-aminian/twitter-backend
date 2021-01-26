package handler

import (
	"fmt"
	"github.com/arman-aminian/twitter-backend/model"
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
	fmt.Println("hahahahahah ", trends)
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

func (h *Handler) GetHashtagTweets(c echo.Context) error {
	// name :=
	return nil
}
