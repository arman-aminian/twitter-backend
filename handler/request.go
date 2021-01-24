package handler

import (
	"github.com/arman-aminian/twitter-backend/model"
	"github.com/labstack/echo/v4"
)

type userRegisterRequest struct {
	User struct {
		Username string `json:"username" bson:"_id" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	} `json:"user"`
}

func (r *userRegisterRequest) bind(c echo.Context, u *model.User) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	print(r.User.Username)
	u.Username = r.User.Username
	u.Email = r.User.Email
	h, err := u.HashPassword(r.User.Password)
	if err != nil {
		return err
	}
	u.Password = h
	return nil
}

type userLoginRequest struct {
	User struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	} `json:"user"`
}

func (r *userLoginRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	return nil
}

type tweetCreateRequest struct {
	Tweet struct {
		Text  string `json:"text" bson:"text"`
		Media string `json:"media" bson:"media"`
	} `json:"tweet"`
}

func (r *tweetCreateRequest) bind(c echo.Context, a *model.Tweet) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	a.Text = r.Tweet.Text
	a.Media = r.Tweet.Media
	return nil
}
