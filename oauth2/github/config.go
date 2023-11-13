package github

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/adoublef/yap/env"
	o2 "github.com/adoublef/yap/oauth2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	api = "https://api.github.com/user"
)

type Config struct {
	oauth2.Config
}

func NewConfig(url string) (c *Config) {
	var (
		client = env.Must("GITHUB_CLIENT_ID")
		secret = env.Must("GITHUB_CLIENT_SECRET")
	)

	c = &Config{
		Config: oauth2.Config{
			ClientID:     client,
			ClientSecret: secret,
			Endpoint:     github.Endpoint,
			Scopes:       []string{"user:email"},
			RedirectURL:  url,
		},
	}
	return
}

func (c *Config) UserInfo(ctx context.Context, tok *oauth2.Token) (user *o2.UserInfo, err error) {
	r, err := c.Client(ctx, tok).Get(api)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	var v User
	if err = json.NewDecoder(r.Body).Decode(&v); err != nil {
		return nil, err
	}
	// TODO return userinfo
	user = &o2.UserInfo{
		ID: o2.ID{
			Provider: "github",
			User:     strconv.Itoa(v.ID)},
		Photo: v.AvatarUrl,
		Login: v.Login,
		Name:  v.Name,
		Email: v.Email,
	}
	return
}

// User is a provider agnostic representation of a user
//
// See (https://github.com/google/go-github/blob/master/github/users.go)
type User struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	AvatarUrl string `json:"avatar_url"`
	Name      string `json:"name"`
	Email     string `json:"email"`
}
