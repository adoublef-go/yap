package http

import (
	"fmt"
	"net/http"

	"github.com/rs/xid"
)

func (s *Service) handleIndex(redirect string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userID xid.ID
		ok := s.us.Has(w, r, &userID)
		fmt.Printf("does a user exist? %v\n\n", ok)
		http.Redirect(w, r, redirect, http.StatusFound)
	}
}
