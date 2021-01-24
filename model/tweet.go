package model

type Tweet struct {
	Text     string    `json:"text" bson:"text"`
	Media    string    `json:"media" bson:"media"`
	Owner    *User     `json:"owner" bson:"owner"`
	Likes    *[]string `json:"likes" bson:"likes"`
	Retweets *[]string `json:"retweets" bson:"retweets"`
}

func NewTweet() *Tweet {
	var t Tweet
	t.Owner = &User{}
	t.Likes = &[]string{}
	t.Retweets = &[]string{}
	return &t
}
