package model

type Owner struct {
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
}

func NewOwner(username, pic string) *Owner {
	return &Owner{
		Username:       username,
		ProfilePicture: pic,
	}
}
