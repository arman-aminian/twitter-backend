package handler

import (
	"encoding/json"
	"github.com/arman-aminian/twitter-backend/db"
	"github.com/arman-aminian/twitter-backend/model"
	"github.com/arman-aminian/twitter-backend/router"
	"github.com/arman-aminian/twitter-backend/store"
	"github.com/arman-aminian/twitter-backend/user"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	usersDb *mongo.Collection
	us      user.Store
	h       *Handler
	e       *echo.Echo
)

func authHeader(token string) string {
	return "Token " + token
}

func setup(usernames []string, shouldLoadFixtures bool) {
	usersDb = db.SetupUsersDb()
	us = store.NewUserStore(usersDb)
	h = NewHandler(us)
	e = router.New()
	_ = cleanUp(usernames)
	if shouldLoadFixtures {
		_ = loadFixtures()
	}
}

func responseMap(b []byte, key string) map[string]interface{} {
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	return m[key].(map[string]interface{})
}

func loadFixtures() error {
	u1 := model.NewUser()
	u1.Username = "user1"
	u1.Email = "user1@realworld.io"
	u1.Bio = "user1 bio"
	u1.ProfilePicture = "https://aux.iconspalace.com/uploads/18923702171865348111.png"
	u1.HeaderPicture = "https://www.polystar.com/wp-content/uploads/2019/01/polystar-solutions-header.jpg"
	u1.Password, _ = u1.HashPassword("user1_pass")
	if err := us.Create(u1); err != nil {
		return err
	}

	u2 := model.NewUser()
	u2.Username = "user2"
	u2.Email = "user2@realworld.io"
	u2.Bio = "user2 bio"
	u2.ProfilePicture = "https://cdn4.iconfinder.com/data/icons/small-n-flat/24/user-alt-512.png"
	u2.HeaderPicture = "https://blog.hubspot.com/hubfs/best-twitter-cover-photo-headers.jpg"
	u2.Password, _ = u2.HashPassword("user2_pass")
	if err := us.Create(u2); err != nil {
		return err
	}

	err := us.AddFollower(u2, u1)

	return err
}

func cleanUp(usernames []string) error {
	for _, username := range usernames {
		err := us.Remove(username)
		if err != nil {
			return err
		}
	}
	return nil
}
