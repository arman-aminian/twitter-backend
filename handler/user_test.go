package handler

import (
	"fmt"
	"github.com/arman-aminian/twitter-backend/router/middleware"
	"github.com/arman-aminian/twitter-backend/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func helper(usernames []string, shouldLoadFixtures bool, reqJSON, target string) (echo.Context, *httptest.ResponseRecorder) {
	setup(usernames, shouldLoadFixtures)
	req := httptest.NewRequest(echo.POST, target, strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

func TestSignUpSuccess(t *testing.T) {
	username := "alice"
	reqJSON := fmt.Sprintf(`{"user":{"username":"%s","email":"%s@realworld.io","password":"%s_pass"}}`, username, username, username)
	c, rec := helper([]string{username}, false, reqJSON, "/signup")
	assert.NoError(t, h.SignUp(c))
	if assert.Equal(t, http.StatusCreated, rec.Code) {
		m := responseMap(rec.Body.Bytes(), "user")
		assert.Equal(t, fmt.Sprintf("%s", username), m["username"])
		assert.Equal(t, fmt.Sprintf("%s@realworld.io", username), m["email"])
		assert.Nil(t, m["bio"])
		assert.Nil(t, m["profile_picture"])
		assert.Nil(t, m["header_picture"])
		assert.Nil(t, m["tweets"])
		assert.Nil(t, m["followings"])
		assert.Nil(t, m["followers"])
		assert.Nil(t, m["notifications"])
		assert.Nil(t, m["logs"])
		assert.NotEmpty(t, m["token"])
	}
}

func TestSignUpFailed(t *testing.T) {
	username := "alice"
	reqJSON := fmt.Sprintf(`{"user":{"username":"%s","email":"%s@realworld.io","password":"%s_pass"}}`, username, username, username)
	c, rec := helper(nil, false, reqJSON, "/signup")
	assert.NoError(t, h.SignUp(c))
	assert.Equal(t, rec.Code, http.StatusUnprocessableEntity)
}

func TestLoginSuccess(t *testing.T) {
	username := "alice"
	reqJSON := fmt.Sprintf(`{"user":{"email":"%s@realworld.io","password":"%s_pass"}}`, username, username)
	c, rec := helper(nil, false, reqJSON, "/login")
	assert.NoError(t, h.Login(c))
	if assert.Equal(t, http.StatusOK, rec.Code) {
		m := responseMap(rec.Body.Bytes(), "user")
		assert.Equal(t, fmt.Sprintf("%s", username), m["username"])
		assert.Equal(t, fmt.Sprintf("%s@realworld.io", username), m["email"])
		assert.NotEmpty(t, m["token"])
	}
}

func TestLoginFailedNotFound(t *testing.T) {
	username := "asghar"
	reqJSON := fmt.Sprintf(`{"user":{"email":"%s@realworld.io","password":"%s_pass"}}`, username, username)
	c, rec := helper(nil, false, reqJSON, "/login")
	assert.NoError(t, h.Login(c))
	// 500 is the mongoDB not found error. TODO fix this
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestLoginFailedWrongPassword(t *testing.T) {
	username := "alice"
	reqJSON := fmt.Sprintf(`{"user":{"email":"%s@realworld.io","password":"something_else"}}`, username)
	c, rec := helper(nil, false, reqJSON, "/login")
	assert.NoError(t, h.Login(c))
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestGetProfileSuccess(t *testing.T) {
	setup([]string{"user1", "user2"}, true)
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.GET, "/profiles/:username", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT("user1")))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/profiles/:username")
	c.SetParamNames("username")
	c.SetParamValues("user1")
	err := jwtMiddleware(func(context echo.Context) error {
		return h.GetProfile(c)
	})(c)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusOK, rec.Code) {
		m := responseMap(rec.Body.Bytes(), "profile")
		assert.Equal(t, "user1", m["username"])
		assert.Equal(t, "user1 bio", m["bio"])
		assert.Equal(t, "https://aux.iconspalace.com/uploads/18923702171865348111.png", m["profile_picture"])
		assert.Equal(t, "https://www.polystar.com/wp-content/uploads/2019/01/polystar-solutions-header.jpg", m["header_picture"])
		assert.Empty(t, m["followers"])
		assert.NotEmpty(t, m["followings"])
		assert.Equal(t, false, m["is_following"])
	}
	// _ = cleanUp([]string{"user1", "user2"})
}
