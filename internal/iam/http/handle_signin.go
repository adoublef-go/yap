package http

import (
	"net/http"

	"github.com/adoublef/yap/oauth2"
)

func (s *Service) handleSignIn(c oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := s.a.SignIn(w, r, c)
		if err != nil {
			http.Error(w, ErrAuthCodeURL, http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, url, http.StatusFound)
	}
}
