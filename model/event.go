package model

type Event struct {
	Mode    string `json:"mode" bson:"mode"`
	Partner string `json:"partner" bson:"partner"`
	Content string `json:"content" bson:"content"`
	EventID string `json:"event_id" bson:"event_id"`
}
