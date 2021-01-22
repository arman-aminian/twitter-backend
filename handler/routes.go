package handler

import (
	"github.com/labstack/echo/v4"
)

const (
	signUp = "/signup"
	login  = "/login"
)

func (h *Handler) Register(g *echo.Group) {
	g.POST(signUp, h.SignUp)
	g.POST(login, h.Login)
}
