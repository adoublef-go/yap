package iam

import (
	"github.com/adoublef/yap/oauth2"
	"github.com/rs/xid"
)

type User struct {
	Profile *Profile  `json:"profile"`
	Email   string    `json:"email"`
	OAuth2  oauth2.ID `json:"oauth2"`
}

func UserFromOAuth2(info *oauth2.UserInfo) *User {
	u := &User{
		Profile: &Profile{
			ID:    xid.New(),
			Photo: info.Photo,
			Name:  info.Name,
			Login: info.Login,
		},
		Email:  info.Email,
		OAuth2: info.ID,
	}
	return u
}

func (u User) ID() xid.ID {
	return u.Profile.ID
}

type Profile struct {
	ID    xid.ID `json:"id"`
	Login string `json:"login"`
	Photo string `json:"photoUrl"`
	Name  string `json:"name"`
}
