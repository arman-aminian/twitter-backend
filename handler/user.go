package handler

import (
	"errors"
	"github.com/arman-aminian/twitter-backend/model"
	"github.com/arman-aminian/twitter-backend/utils"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
)

func (h *Handler) Dummy(c echo.Context) error {
	return c.JSON(http.StatusCreated, errors.New("hello world"))
}

// signUp godoc
// @Summary Register a new user
// @Description Register a new user
// @ID sign-up
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body userRegisterRequest true "User info for registration: email, username and password"
// @Success 201 {object} userResponse
// @Failure 404 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Router /signup [post]
func (h *Handler) SignUp(c echo.Context) error {
	u := model.NewUser()
	req := &userRegisterRequest{}
	if err := req.bind(c, u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if err := h.userStore.Create(u); err != nil {
		return c.JSON(http.StatusNotFound, utils.NewError(err))
	}
	response := newUserResponse(u)
	//cookie := new(http.Cookie)
	//cookie.Name = "Token"
	//cookie.Value = response.User.Token
	//cookie.Expires = time.Now().Add(24 * time.Hour)
	//c.SetCookie(cookie)

	//header('Access-Control-Allow-Origin', yourExactHostname);

	//c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "http://localhost:3000")
	//c.Response().Header().Add(echo.HeaderAccessControlAllowCredentials, "true")
	//c.Response().Header().Add(echo.HeaderAccessControlAllowOrigin, "http://localhost:3000")
	//c.Response().Header().Add(echo.HeaderAccessControlAllowHeaders, "Origin, X-Requested-With, Content-Type, Accept")
	//c.Response().Header().
	return c.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary Login an existing user
// @Description Login an existing user
// @ID login
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body userLoginRequest true "Credentials to use: email and password"
// @Success 200 {object} userResponse
// @Failure 403 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /login [post]
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
	response := newUserResponse(u)
	//cookie := new(http.Cookie)
	//cookie.Name = "Token"
	//cookie.Value = response.User.Token
	//cookie.Expires = time.Now().Add(24 * time.Hour)
	//c.SetCookie(cookie)
	return c.JSON(http.StatusCreated, response)
}

// UpdateUser godoc
// @Summary Update User
// @Description Update user information for a user whom information is given
// @ID update-user
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body userUpdateRequest true "User details to update. At least **one** field is required."
// @Success 200 {object} userResponse
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

	newUser.Username = oldUser.Username
	newUser.Tweets = oldUser.Tweets
	newUser.Followings = oldUser.Followings
	newUser.Followers = oldUser.Followers
	newUser.Bio = oldUser.Bio
	newUser.Name = oldUser.Name
	newUser.ProfilePicture = oldUser.ProfilePicture
	newUser.HeaderPicture = oldUser.HeaderPicture
	newUser.Password = oldUser.Password
	newUser.Email = oldUser.Email
	newUser.Logs = oldUser.Logs
	newUser.Notifications = oldUser.Notifications

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
// @Summary Get a User profile
// @Description Get a profile of a user of the system. Authorization is optional
// @ID get-profile
// @Tags profile
// @Accept  json
// @Produce  json
// @Success 200 {object} userResponse
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
// @Description Update user profile based on 4 entries: name, bio, and pictures
// @ID update-profile
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body userProfileUpdateRequest true "User details to update. At least **one** field is required."
// @Success 200 {object} userResponse
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /profiles/{username} [put]
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

// GetProfilePictureFile godoc
// @Summary Get the profile picture
// @Description Get the profile picture of a user based on the filename
// @ID get-profile-picture
// @Tags user
// @Param user body string true "name of the file."
// @Router /profile-pictures/{filename} [get]
func (h *Handler) GetProfilePictureFile(c echo.Context) error {
	mediaFolderName := "media/profile-pictures/"
	mediaPath := mediaFolderName + c.Param("filename")
	return c.File(mediaPath)
}

// GetHeaderPictureFile godoc
// @Summary Get the header picture
// @Description Get the header picture of a user based on the filename
// @ID get-header-picture
// @Tags user
// @Param user body string true "name of the file."
// @Router /header-pictures/{filename} [get]
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
// @Success 201 {object} profileResponse
// @Failure 400 {object} utils.Error
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

// GetTimeline godoc
// @Summary Get the timeline of your account
// @Description Get tweets of your followings. Auth is required.
// @ID get-timeline
// @Tags timeline
// @Produce  json
// @Success 200 {object} tweetListResponse
// @Failure 422 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /home [get]
func (h *Handler) GetTimeline(c echo.Context) error {
	day, err := strconv.Atoi(c.Param("day"))
	if err != nil {
		day = 0
	}
	day = -1 * day
	u, err := h.userStore.GetByUsername(stringFieldFromToken(c, "username"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	var usernames []string
	for _, f := range *u.Followings {
		usernames = append(usernames, f.Username)
	}
	usernames = append(usernames, u.Username)
	if len(usernames) == 0 {
		return c.JSON(http.StatusOK, newTweetListResponse(stringFieldFromToken(c, "username"), nil, 0))
	}

	tweetsId, err := h.userStore.GetTweetIdListFromUsernameList(usernames)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if len(*tweetsId) == 0 {
		return c.JSON(http.StatusOK, newTweetListResponse(stringFieldFromToken(c, "username"), nil, 0))
	}

	timelineTweets, err := h.tweetStore.GetTimelineFromTweetIDs(*tweetsId, day)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	// sort timeline by tweet's creating time
	timeline := *timelineTweets
	sort.Slice(timeline, func(i, j int) bool {
		return timeline[i].Time.After(timeline[j].Time)
	})
	return c.JSON(http.StatusOK, newTweetListResponse(stringFieldFromToken(c, "username"), &timeline, len(timeline)))
}

// SearchUsernames godoc
// @Summary Search a user by its username
// @Description Search by specifying the username. Auth is optional.
// @ID search-username
// @Tags search
// @Produce  json
// @Success 200 {object} tweetListResponse
// @Failure 400 {object} utils.Error
// @Router /search/username [get]
func (h *Handler) SearchUsernames(c echo.Context) error {
	query := c.QueryParam("query")
	if query == "" {
		return c.JSON(http.StatusBadRequest, utils.NewError(errors.New("nothing to search for")))
	}
	result, err := h.userStore.GetUsernameSearchResult(query)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, newOwnerList(h.userStore, stringFieldFromToken(c, "username"), result))
}

// SearchTweets godoc
// @Summary Search tweets related to query
// @Description Search all the tweets related to the query.
// @ID search-tweets
// @Tags search
// @Accept json
// @Produce  json
// @Param query query string true "Part of the tweets you want to search."
// @Success 200 {object} tweetListResponse
// @Failure 400 {object} utils.Error
// @Router /search/tweet [get]
func (h *Handler) SearchTweets(c echo.Context) error {
	query := &model.SearchQuery{}
	err := c.Bind(query)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if query.Query == "" {
		return c.JSON(http.StatusBadRequest, utils.NewError(errors.New("nothing to search for")))
	}
	result, err := h.tweetStore.GetTweetSearchResult(query.Query)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, newTweetListResponse(stringFieldFromToken(c, "username"), result, len(*result)))
}

// GetFollowingAndFollowersList godoc
// @Summary Get the list of followings and followers of a user. Auth is required.
// @Description Get the list of followings and followers of a user. Auth is required.
// @ID get-following-followers
// @Tags following-followers
// @Produce  json
// @Success 200 {object} FollowersAndFollowingListResponse
// @Failure 400 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Router /tweets/{id}/list [get]
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

// GetLogs godoc
// @Description Get the list of everything the user has done (logs).
// @ID get-logs
// @Tags logs
// @Produce  json
// @Success 200 {object} EventListResponse
// @Failure 400 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /profiles/{username}/logs [get]
func (h *Handler) GetLogs(c echo.Context) error {
	u, err := h.userStore.GetByUsername(c.Param("username"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	if u.Username != stringFieldFromToken(c, "username") {
		return c.JSON(http.StatusBadRequest, errors.New("can't get another user's logs"))
	}
	return c.JSON(http.StatusOK, newLogsList(u))
}

// GetNotifications godoc
// @Description Get the list of everything done with the use being the target (notifications).
// @ID get-notifications
// @Tags notifications
// @Produce  json
// @Success 200 {object} EventListResponse
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /profiles/{username}/notifications [get]
func (h *Handler) GetNotifications(c echo.Context) error {
	u, err := h.userStore.GetByUsername(c.Param("username"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	return c.JSON(http.StatusOK, newNotificationsList(u))
}

// GetSuggestions godoc
// @Description Get the list of the users suggested to a user (followers of its followings)
// @ID get-suggestions
// @Tags suggestions
// @Produce  json
// @Success 200 {object} OwnerListResponse
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /suggestions [get]
func (h *Handler) GetSuggestions(c echo.Context) error {
	username := stringFieldFromToken(c, "username")
	u, err := h.userStore.GetByUsername(username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	if len(*u.Followings) == 0 {
		return c.JSON(http.StatusOK, newOwnerList(h.userStore, username, nil))
	}
	var suggestions []model.Owner
	followings := *u.Followings
	for _, f := range followings {
		following, _ := h.userStore.GetByUsername(f.Username)
		suggestions = append(suggestions, *following.Followings...)
	}

	suggestionsSize := len(suggestions)
	followingsSize := len(followings)
	for i := range suggestions {
		for j := range followings {
			indexI := suggestionsSize - 1 - i
			indexJ := followingsSize - 1 - j
			if suggestions[indexI].Username == followings[indexJ].Username || suggestions[indexI].Username == u.Username {
				suggestions = removeIndex(suggestions, indexI)
				break
			}
		}
	}

	// to sort suggestions by their frequencies
	suggestionsFreq := dupCount(suggestions)
	sorted := make([]model.Owner, 0, len(suggestionsFreq))
	for name := range suggestionsFreq {
		sorted = append(sorted, name)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return suggestionsFreq[sorted[i]] > suggestionsFreq[sorted[j]]
	})

	maxNumberOfSuggestions := 3
	if len(sorted) < maxNumberOfSuggestions {
		return c.JSON(http.StatusOK, newOwnerList(h.userStore, username, &sorted))
	}
	sorted = sorted[:maxNumberOfSuggestions]
	return c.JSON(http.StatusOK, newOwnerList(h.userStore, username, &sorted))
}

func removeIndex(s []model.Owner, index int) []model.Owner {
	return append(s[:index], s[index+1:]...)
}

func dupCount(list []model.Owner) map[model.Owner]int {
	duplicateFrequency := make(map[model.Owner]int)
	for _, item := range list {
		_, exist := duplicateFrequency[item]
		if exist {
			duplicateFrequency[item] += 1 // increase counter by 1 if already in the map
		} else {
			duplicateFrequency[item] = 1 // else start counting from 1
		}
	}
	return duplicateFrequency
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
