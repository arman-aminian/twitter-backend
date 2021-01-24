package handler

import (
	"github.com/arman-aminian/twitter-backend/model"
	"github.com/labstack/echo/v4"
)

// Registration request
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

// Login request
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

// Profile update request
type userProfileUpdateRequest struct {
	User struct {
		Bio            string `json:"bio"`
		ProfilePicture string `json:"profile_picture"`
		HeaderPicture  string `json:"header_Picture"`
	} `json:"user"`
}

func newUserProfileUpdateRequest() *userProfileUpdateRequest {
	return new(userProfileUpdateRequest)
}

func (r *userProfileUpdateRequest) populate(u *model.User) {
	if len(u.Bio) <= 160 {
		r.User.Bio = u.Bio
	}
	if len(u.ProfilePicture) > 0 {
		r.User.ProfilePicture = u.ProfilePicture
	}
	if len(u.HeaderPicture) > 0 {
		r.User.HeaderPicture = u.HeaderPicture
	}
}

func (r *userProfileUpdateRequest) bind(c echo.Context, u *model.User) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	u.Bio = r.User.Bio
	u.ProfilePicture = r.User.ProfilePicture
	u.HeaderPicture = r.User.HeaderPicture
	return nil
}

// User update request
type userUpdateRequest struct {
	User struct {
		Username string `json:"username"`
		Email    string `json:"email" validate:"email"`
		Password string `json:"password"`
	} `json:"user"`
}

func newUserUpdateRequest() *userUpdateRequest {
	return new(userUpdateRequest)
}

func (r *userUpdateRequest) populate(u *model.User) {
	r.User.Username = u.Username
	r.User.Email = u.Email
	r.User.Password = u.Password
}

func (r *userUpdateRequest) bind(c echo.Context, u *model.User) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	u.Username = r.User.Username
	u.Email = r.User.Email
	if r.User.Password != u.Password {
		h, err := u.HashPassword(r.User.Password)
		if err != nil {
			return err
		}
		u.Password = h
	}
	return nil
}

// Tweet request section start from here

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
