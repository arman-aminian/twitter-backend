package model

import (
	"time"
)

type Event struct {
	Mode      string    `json:"mode" bson:"mode"`
	Source    Owner     `json:"source" bson:"source"`
	Target    Owner     `json:"target" bson:"target"`
	Content   string    `json:"content" bson:"content"`
	TimeStamp time.Time `json:"timestamp" bson:"timestamp"`
	Tweet     *Tweet    `json:"tweet" bson:"tweet"`
}

func NewEvent() *Event {
	var e Event
	e.TimeStamp, _ = time.Parse("2006-01-02 03:04:05", time.Now().Format("2006-01-02 03:04:05"))
	e.Tweet = &Tweet{}
	return &e
}
