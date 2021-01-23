package handler

import (
	"github.com/arman-aminian/twitter-backend/router/middleware"
	"github.com/arman-aminian/twitter-backend/utils"
	"github.com/labstack/echo/v4"
)

const (
	signUp   = "/signup"
	login    = "/login"
	profiles = "/profiles"
)

func (h *Handler) Register(g *echo.Group) {
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	g.POST(signUp, h.SignUp)
	g.POST(login, h.Login)

	profiles := g.Group(profiles, jwtMiddleware)
	profiles.GET("/:username", h.GetProfile)
}
