package handler

import (
	"github.com/arman-aminian/twitter-backend/router/middleware"
	"github.com/arman-aminian/twitter-backend/utils"
	"github.com/labstack/echo/v4"
)

const (
	signUp    = "/signup"
	login     = "/login"
	timeline  = "/home"
	userPath  = "/user"
	profiles  = "/profiles"
	tweets    = "/tweets"
	usernameQ = "/:username"
	follow    = usernameQ + "/follow"
	media     = "/media"
)

func (h *Handler) Register(g *echo.Group) {
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	g.POST(signUp, h.SignUp)
	g.POST(login, h.Login)

	home := g.Group(timeline, jwtMiddleware)
	home.GET("", h.GetTimeline)

	user := g.Group(userPath, jwtMiddleware)
	user.PUT(usernameQ, h.UpdateUser)

	profiles := g.Group(profiles, jwtMiddleware)
	profiles.GET(usernameQ, h.GetProfile)
	profiles.PUT(usernameQ, h.UpdateProfile)
	profiles.POST(follow, h.Follow)
	profiles.DELETE(follow, h.UnFollow)
	profiles.GET(usernameQ+"/list", h.GetFollowingAndFollowersList)

	tweets := g.Group(tweets, middleware.JWTWithConfig(
		middleware.JWTConfig{
			Skipper: func(c echo.Context) bool {
				// TODO replace INJA and uncomment
				// if c.Request().Method == "GET" && c.Path() != "/tweets/INJA" {
				//	return true
				// }
				return false
			},
			SigningKey: utils.JWTSecret,
		},
	))
	tweets.POST("", h.CreateTweet)
	tweets.GET("/:id", h.GetTweet)
	tweets.GET("/:id/list", h.GetTweetLikeAndRetweetList)
	tweets.POST("/:id/like", h.Like)
	tweets.DELETE("/:id/like", h.UnLike)
	tweets.POST("/:id/retweet", h.Retweet)
	tweets.DELETE("/:id/retweet", h.UnRetweet)

	files := g.Group(media)
	files.GET("/:filename", h.GetFile)

}
