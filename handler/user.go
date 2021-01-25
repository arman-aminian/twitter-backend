package handler

import (
	"errors"
	"github.com/arman-aminian/twitter-backend/model"
	"github.com/arman-aminian/twitter-backend/utils"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
)

// signUp godoc
// @Summary Register a new user
// @Description Register a new user
// @ID sign-up
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body userRegisterRequest true "User info for registration"
// @Success 201 {object} userResponse
// @Failure 400 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /users [post]
func (h *Handler) SignUp(c echo.Context) error {
	u := model.NewUser()
	req := &userRegisterRequest{}
	if err := req.bind(c, u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if err := h.userStore.Create(u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.JSON(http.StatusCreated, newUserResponse(u))
}

// Login godoc
// @Summary Login for existing user
// @Description Login for existing user
// @ID login
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body userLoginRequest true "Credentials to use"
// @Success 200 {object} userResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /users/login [post]
func (h *Handler) Login(c echo.Context) error {
	req := &userLoginRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	u, err := h.userStore.GetByEmail(req.User.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusForbidden, utils.AccessForbidden())
	}
	if !u.CheckPassword(req.User.Password) {
		return c.JSON(http.StatusForbidden, utils.AccessForbidden())
	}
	return c.JSON(http.StatusOK, newUserResponse(u))
}

// UpdateUser godoc
// @Summary Update current user
// @Description Update user information for current user
// @ID update-user
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body userUpdateRequest true "User details to update. At least **one** field is required."
// @Success 200 {object} userResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /user [put]
func (h *Handler) UpdateUser(c echo.Context) error {
	oldUser, err := h.userStore.GetByUsername(c.Param("username"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if oldUser == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	newUser := model.NewUser()
	_ = copier.Copy(&newUser, &oldUser)
	req := newUserUpdateRequest()
	req.populate(newUser)
	if err := req.bind(c, newUser); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if err := h.userStore.Update(oldUser, newUser); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, newUserResponse(newUser))
}

// GetProfile godoc
// @Summary Get a profile
// @Description Get a profile of a user of the system. Auth is optional
// @ID get-profile
// @Tags profile
// @Accept  json
// @Produce  json
// @Param username path string true "Username of the profile to get"
// @Success 200 {object} userResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /profiles/{username} [get]
func (h *Handler) GetProfile(c echo.Context) error {
	destUsername := c.Param("username")
	u, err := h.userStore.GetByUsername(destUsername)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	return c.JSON(http.StatusOK, newProfileResponse(h.userStore, stringFieldFromToken(c, "username"), u))
}

// UpdateProfile godoc
// @Summary Update a user's profile
// @Description Update user profile
// @ID update-profile
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body userProfileUpdateRequest true "User details to update. At least **one** field is required."
// @Success 200 {object} userResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /user [put]
func (h *Handler) UpdateProfile(c echo.Context) error {
	u, err := h.userStore.GetByUsername(stringFieldFromToken(c, "username"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	req := newUserProfileUpdateRequest()
	req.populate(u)
	u.Name = c.FormValue("name")
	u.Bio = c.FormValue("bio")
	ppf, err := c.FormFile("profile_picture")
	if err == nil {
		src, err := ppf.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err))
		}
		defer src.Close()

		mediaFolderName := "media/profile-pictures/"
		mediaPath := mediaFolderName + ppf.Filename
		dst, err := os.Create(mediaPath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err))
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err))
		}
		u.ProfilePicture = mediaPath
	} else {
		// Update without Profile Picture
		u.ProfilePicture = ""
	}

	hpf, err := c.FormFile("header_picture")
	if err == nil {
		src, err := hpf.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err))
		}
		defer src.Close()

		mediaFolderName := "media/header-pictures/"
		mediaPath := mediaFolderName + hpf.Filename
		dst, err := os.Create(mediaPath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err))
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err))
		}
		u.HeaderPicture = mediaPath
	} else {
		// Update without Header Picture
		u.HeaderPicture = ""
	}

	if err := h.userStore.UpdateProfile(u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, newProfileResponse(h.userStore, u.Username, u))
}

func (h *Handler) GetProfilePictureFile(c echo.Context) error {
	mediaFolderName := "media/profile-pictures/"
	mediaPath := mediaFolderName + c.Param("filename")
	return c.File(mediaPath)
}

func (h *Handler) GetHeaderPictureFile(c echo.Context) error {
	mediaFolderName := "media/header-pictures/"
	mediaPath := mediaFolderName + c.Param("filename")
	return c.File(mediaPath)
}

// Follow godoc
// @Summary Follow a user
// @Description Follow a user by username
// @ID follow
// @Tags follow
// @Accept  json
// @Produce  json
// @Param username path string true "Username of the profile you want to follow"
// @Success 200 {object} profileResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /profiles/{username}/follow [post]
func (h *Handler) Follow(c echo.Context) error {
	follower, err := h.userStore.GetByUsername(stringFieldFromToken(c, "username"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	u, err := h.userStore.GetByUsername(c.Param("username"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	if u.Username == follower.Username {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("can't follow yourself")))
	}
	if Contains(*u.Followers, follower.Username) || Contains(*follower.Followings, u.Username) {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("already follows the target")))
	}

	if err := h.userStore.AddFollower(u, follower); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	e := h.CreateFollowEvent(follower, u)
	err = h.userStore.AddLog(follower, e)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	err = h.userStore.AddNotification(u, e)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, newProfileResponse(h.userStore, follower.Username, u))
}

// Unfollow godoc
// @Summary Unfollow a user
// @Description Unfollow a user by username
// @ID unfollow
// @Tags follow
// @Accept  json
// @Produce  json
// @Param username path string true "Username of the profile you want to unfollow"
// @Success 201 {object} userResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /profiles/{username}/follow [delete]
func (h *Handler) UnFollow(c echo.Context) error {
	follower, err := h.userStore.GetByUsername(stringFieldFromToken(c, "username"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	u, err := h.userStore.GetByUsername(c.Param("username"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	if u.Username == follower.Username {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("can't unfollow yourself")))
	}
	if !Contains(*u.Followers, follower.Username) || !Contains(*follower.Followings, u.Username) {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("doesn't follow the target")))
	}
	if err := h.userStore.RemoveFollower(u, follower); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, newProfileResponse(h.userStore, stringFieldFromToken(c, "username"), u))
}

// Articles godoc
// @Summary Get recent articles globally
// @Description Get most recent articles globally. Use query parameters to filter results. Auth is optional
// @ID get-articles
// @Tags article
// @Accept  json
// @Produce  json
// @Param tag query string false "Filter by tag"
// @Param author query string false "Filter by author (username)"
// @Param favorited query string false "Filter by favorites of a user (username)"
// @Param limit query integer false "Limit number of articles returned (default is 20)"
// @Param offset query integer false "Offset/skip number of articles (default is 0)"
// @Success 200 {object} articleListResponse
// @Failure 500 {object} utils.Error
// @Router /articles [get]
func (h *Handler) GetTimeline(c echo.Context) error {

	u, err := h.userStore.GetByUsername(stringFieldFromToken(c, "username"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	var usernames []string
	for _, f := range *u.Followings {
		usernames = append(usernames, f.Username)
	}
	if len(usernames) == 0 {
		return c.JSON(http.StatusOK, newTweetListResponse(c, stringFieldFromToken(c, "username"), nil, 0))
	}
	tweetsId, err := h.userStore.GetTweetIdListFromUsernameList(usernames)
	if err != nil {
		return err
	}
	tweets, err := h.tweetStore.GetTimelineFromUsernames(*tweetsId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, newTweetListResponse(c, stringFieldFromToken(c, "username"), tweets, len(*tweets)))
}

func (h *Handler) GetFollowingAndFollowersList(c echo.Context) error {
	u, err := h.userStore.GetByUsername(c.Param("username"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	return c.JSON(http.StatusOK, newFollowingAndFollowersList(h.userStore, stringFieldFromToken(c, "username"), u))
}

func stringFieldFromToken(c echo.Context, field string) string {
	field, ok := c.Get(field).(string)
	if !ok {
		return ""
	}
	return field
}

func Contains(slice []model.Owner, val string) bool {
	for _, item := range slice {
		if item.Username == val {
			return true
		}
	}
	return false
}
