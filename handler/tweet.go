package handler

import (
	"errors"
	"fmt"
	"github.com/arman-aminian/twitter-backend/model"
	"github.com/arman-aminian/twitter-backend/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

// CreateArticle godoc
// @Summary Create an tweet
// @Description Create an tweet
// @ID create-tweet
// @Tags article
// @Accept  json
// @Produce  json
// @Param article body tweetCreateRequest true "Article to create"
// @Success 201 {object} singleTweetResponse
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /articles [post]
func (h *Handler) CreateTweet(c echo.Context) error {
	t := model.NewTweet()

	req := &tweetCreateRequest{}
	if err := req.bind(c, t); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	u, _ := h.userStore.GetByUsername(stringFieldFromToken(c, "username"))
	t.Owner.Username = u.Username
	t.Owner.ProfilePicture = u.ProfilePicture
	t.Time = time.Now()
	t.Date = time.Now().Format("2006-01-02")
	t.ID = primitive.NewObjectID()
	err := h.tweetStore.CreateTweet(t)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	//print(a.OwnerUsername)
	err = h.userStore.AddTweet(u, t)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	res := newTweetResponse(c, t)
	fmt.Println(res)
	return c.JSON(http.StatusCreated, res)
}

// GetArticle godoc
// @Summary Get an article
// @Description Get an article. Auth not required
// @ID get-article
// @Tags article
// @Accept  json
// @Produce  json
// @Param slug path string true "Slug of the article to get"
// @Success 200 {object} singleArticleResponse
// @Failure 400 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /articles/{slug} [get]
func (h *Handler) GetTweetLikeAndRetweetList(c echo.Context) error {
	id := c.Param("id")
	t, err := h.tweetStore.GetTweetById(&id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if t == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	return c.JSON(http.StatusOK, newLikeAndRetweetResponse(t))
}

// GetArticle godoc
// @Summary Get an article
// @Description Get an article. Auth not required
// @ID get-article
// @Tags article
// @Accept  json
// @Produce  json
// @Param slug path string true "Slug of the article to get"
// @Success 200 {object} singleArticleResponse
// @Failure 400 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /articles/{slug} [get]
func (h *Handler) GetRetweetList(c echo.Context) error {
	id := c.Param("id")
	t, err := h.tweetStore.GetTweetById(&id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if t == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	return c.JSON(http.StatusOK, newLikeAndRetweetResponse(t))
}

// Favorite godoc
// @Summary Favorite an article
// @Description Favorite an article. Auth is required
// @ID favorite
// @Tags favorite
// @Accept  json
// @Produce  json
// @Param slug path string true "Slug of the article that you want to favorite"
// @Success 200 {object} singleArticleResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /articles/{slug}/favorite [post]
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

	return c.JSON(http.StatusOK, newTweetResponse(c, t))
}

// Favorite godoc
// @Summary Favorite an article
// @Description Favorite an article. Auth is required
// @ID favorite
// @Tags favorite
// @Accept  json
// @Produce  json
// @Param slug path string true "Slug of the article that you want to favorite"
// @Success 200 {object} singleArticleResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /articles/{slug}/favorite [post]
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

// Favorite godoc
// @Summary Favorite an article
// @Description Favorite an article. Auth is required
// @ID favorite
// @Tags favorite
// @Accept  json
// @Produce  json
// @Param slug path string true "Slug of the article that you want to favorite"
// @Success 200 {object} singleArticleResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /articles/{slug}/favorite [post]
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

	return c.JSON(http.StatusOK, newTweetResponse(c, t))
}

// Favorite godoc
// @Summary Favorite an article
// @Description Favorite an article. Auth is required
// @ID favorite
// @Tags favorite
// @Accept  json
// @Produce  json
// @Param slug path string true "Slug of the article that you want to favorite"
// @Success 200 {object} singleArticleResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /articles/{slug}/favorite [post]
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

	return c.JSON(http.StatusOK, newTweetResponse(c, t))
}
