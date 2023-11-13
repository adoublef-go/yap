package http

import (
	"net/http"

	cookie "github.com/adoublef/yap/http"
	"github.com/adoublef/yap/oauth2"
	"github.com/adoublef/yap/oauth2/github"
	sql "github.com/adoublef/yap/sqlite3"
	"github.com/go-chi/chi/v5"
)

var (
	ErrAuthCodeURL    = "Failed to create auth code url"
	ErrAuthentication = "Failed to authenticate"
	ErrDatabaseQuery  = "Failed to complete query"
)

type Service struct {
	m  *chi.Mux
	a  *oauth2.Authenticator
	db sql.ReadWriter
	// us handles user-session cookies
	us *cookie.CookieJar
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.m.ServeHTTP(w, r)
}

func NewService(db sql.ReadWriter, jar *cookie.CookieJar) (s *Service) {
	s = &Service{
		m:  chi.NewMux(),
		db: db,
		us: jar,
		a:  &oauth2.Authenticator{},
	}
	s.routes()
	return
}

func (s *Service) routes() {
	// if logged in, redirect
	s.m.Get("/", s.handleIndex("/rooms"))
	// signin | callback | signout
	// replay onto primary(?)
	gg := github.NewConfig("http://localhost:8080/callback/github")
	s.m.Get("/signin/github", s.handleSignIn(gg))
	s.m.Get("/callback/github", s.handleCallback(gg))
}

type SessionKV struct{}
