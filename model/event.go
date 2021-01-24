package model

import (
	"time"
)

type Event struct {
	Mode      string    `json:"mode" bson:"mode"`
	Source    string    `json:"source" bson:"source"`
	Target    string    `json:"target" bson:"target"`
	Content   string    `json:"content" bson:"content"`
	TimeStamp time.Time `json:"timestamp" bson:"timestamp"`
	// EventID   string              `json:"event_id" bson:"event_id"`
}

func NewEvent() *Event {
	var e Event
	e.TimeStamp, _ = time.Parse("2006-01-02 03:04:05", time.Now().Format("2006-01-02 03:04:05"))
	return &e
}
