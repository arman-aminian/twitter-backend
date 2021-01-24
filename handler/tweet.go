package handler

import (
	"fmt"
	"github.com/arman-aminian/twitter-backend/model"
	"github.com/arman-aminian/twitter-backend/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// CreateArticle godoc
// @Summary Create an article
// @Description Create an article. Auth is require
// @ID create-article
// @Tags article
// @Accept  json
// @Produce  json
// @Param article body articleCreateRequest true "Article to create"
// @Success 201 {object} singleArticleResponse
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

	u, _ := h.userStore.GetByUsername(usernameFromToken(c))
	t.Owner.Username = u.Username
	t.Owner.ProfilePicture = u.ProfilePicture
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
