package model

type Owner struct {
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
	Name           string `json:"name"`
	Bio            string `json:"bio"`
}

func NewOwner(username string, pic string, name string, bio string) *Owner {
	return &Owner{
		Username:       username,
		ProfilePicture: pic,
		Name:           name,
		Bio:            bio,
	}
}
