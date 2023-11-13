package oauth2

import (
	"context"
	"crypto/subtle"
	"database/sql/driver"
	"errors"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

type Authenticator struct {
	// try with one config
	// c
}

func (a Authenticator) SignIn(w http.ResponseWriter, r *http.Request, cfg Config) (url string, err error) {
	state := oauth2.GenerateVerifier()
	c := http.Cookie{
		Name:     cookieName(r, oauth),
		Secure:   r.TLS != nil,
		Path:     "/",
		HttpOnly: true,
		Value:    state,
		// 10 minutes is recommended
		MaxAge:   int((10 * time.Minute).Seconds()),
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &c)
	return cfg.AuthCodeURL(state), nil
}

func (a Authenticator) HandleCallback(w http.ResponseWriter, r *http.Request, cfg Config) (u *UserInfo, err error) {
	// get cookie
	cookie, err := getCookie(r, oauth)
	if err != nil {
		return nil, err
	}
	// compare with url of state on request
	if !compare(cookie.Value, r.FormValue("state")) {
		return nil, errors.New("state value mismatch")
	}
	cookie.Value = ""
	cookie.MaxAge = -1
	http.SetCookie(w, cookie) // set cookie
	// Use the custom HTTP client when requesting a token.
	httpClient := &http.Client{Timeout: 2 * time.Second}
	ctx := context.WithValue(r.Context(), oauth2.HTTPClient, httpClient)
	// exchange `code` for `tok`
	tok, err := cfg.Exchange(ctx, r.FormValue("code"))
	if err != nil {
		return nil, err
	}
	// get `userinfo`
	u, err = cfg.UserInfo(ctx, tok)
	return
}

// TODO `signout`

type Config interface {
	UserInfo(ctx context.Context, tok *oauth2.Token) (*UserInfo, error)
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
}

type ID struct {
	Provider string
	User     string
}

func (id ID) Value() (driver.Value, error) {
	return id.String(), nil
}

func (id ID) String() string {
	return id.Provider + ":" + id.User
}

type UserInfo struct {
	ID    ID     `json:"id"`
	Photo string `json:"photo"`
	Login string `json:"login"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func cookieName(r *http.Request, name string) string {
	// NOTE this may not work behind Fly.io
	if secure := r.TLS != nil; secure {
		name = "__Host-" + name
	}
	return name
}

func getCookie(r *http.Request, name string) (*http.Cookie, error) {
	name = cookieName(r, name)
	return r.Cookie(name)
}

func compare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) != 0
}

const (
	oauth = "oauth-session"
)
