package handler

import (
	"fmt"
	"github.com/arman-aminian/twitter-backend/model"
)

func CreateFollowLogEvent(src *model.User, target *model.User) *model.Event {
	e := model.NewEvent()
	e.Mode = "Follow"
	e.Source = *model.NewOwner(src.Username, src.ProfilePicture)
	e.Target = *model.NewOwner(target.Username, target.ProfilePicture)
	e.Content = fmt.Sprintf("User %s followed User %s at %s", e.Source, e.Target, e.TimeStamp)
	return e
}

func CreateFollowNotificationEvent(src *model.User, target *model.User) *model.Event {
	e := model.NewEvent()
	e.Mode = "Follow"
	e.Source = *model.NewOwner(src.Username, src.ProfilePicture)
	e.Target = *model.NewOwner(target.Username, target.ProfilePicture)
	e.Content = fmt.Sprintf("User %s followed you at %s", e.Source, e.TimeStamp)
	return e
}
