package http

import (
	"errors"
	"log"
	"net/http"

	"github.com/adoublef/yap"
	"github.com/adoublef/yap/internal/iam"
	sql "github.com/adoublef/yap/internal/iam/sqlite3"
	"github.com/adoublef/yap/oauth2"
	sqlutil "github.com/adoublef/yap/sqlite3"
)

func (s *Service) handleCallback(c oauth2.Config) http.HandlerFunc {
	// var session = func(w http.ResponseWriter, r *http.Request, u *iam.User) (err error) {
	// 	_, err = s.ss.Put(s.setSession(w, r), u.ID().Bytes())
	// 	return
	// }
	// replay can be moved to middleware if complexity requires it
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx = r.Context()
		)
		v, err := s.a.HandleCallback(w, r, c)
		if err != nil {
			log.Println(err)
			http.Error(w, ErrAuthentication, http.StatusUnauthorized)
			return
		}
		u := iam.UserFromOAuth2(v)
		// NOTE try using an upsert
		switch _, err := sql.LookupProfile(ctx, s.db, u.OAuth2.String()); {
		case errors.Is(err, sqlutil.ErrNoRows):
			if err := sql.RegisterUser(ctx, s.db, u); err != nil {
				http.Error(w, ErrDatabaseQuery, http.StatusInternalServerError)
				return
			}
		case err != nil:
			http.Error(w, ErrDatabaseQuery, http.StatusInternalServerError)
			return
		}
		if err := s.us.Put(w, r, yap.UUID().String(), u.ID()); err != nil {
			http.Error(w, "Failed to set user session", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (s *Service) setSession(w http.ResponseWriter, r *http.Request) (session string) {
	c := http.Cookie{
		Name:     "user-session",
		Value:    yap.UUID().String(),
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &c)
	return c.Value
}

func (s *Service) getSession(w http.ResponseWriter, r *http.Request) (session string) {
	c, err := r.Cookie("site-session")
	if err != nil {
		return ""
	}
	return c.Value
}
