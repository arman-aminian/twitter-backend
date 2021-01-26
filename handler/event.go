package handler

import (
	"fmt"
	"github.com/arman-aminian/twitter-backend/model"
)

func (h *Handler) CreateFollowEvent(src *model.User, target *model.User) *model.Event {
	e := model.NewEvent()
	e.Mode = "Follow"
	e.Source = *model.NewOwner(src.Username, src.ProfilePicture, src.Name, src.Bio)
	e.Target = *model.NewOwner(target.Username, target.ProfilePicture, target.Name, target.Bio)
	e.Content = fmt.Sprintf("User %s followed User %s at %s.", e.Source.Username, e.Target.Username, e.TimeStamp)
	return e
}

func (h *Handler) CreateLikeEvent(src *model.User, t *model.Tweet) *model.Event {
	e := model.NewEvent()
	e.Mode = "Like"
	e.Source = *model.NewOwner(src.Username, src.ProfilePicture, src.Name, src.Bio)
	e.Target = *model.NewOwner(t.Owner.Username, t.Owner.ProfilePicture, t.Owner.Name, t.Owner.Bio)
	e.Content = fmt.Sprintf("User %s liked Tweet %s at %s.", e.Source.Username, t.Text, e.TimeStamp)
	return e
}

func (h *Handler) CreateRetweetEvent(src *model.User, t *model.Tweet) *model.Event {
	e := model.NewEvent()
	e.Mode = "Retweet"
	e.Source = *model.NewOwner(src.Username, src.ProfilePicture, src.Name, src.Bio)
	e.Target = *model.NewOwner(t.Owner.Username, t.Owner.ProfilePicture, t.Owner.Name, t.Owner.Bio)
	e.Content = fmt.Sprintf("User %s retweeted Tweet %s at %s.", e.Source.Username, t.Text, e.TimeStamp)
	return e
}
