package http

import (
	"log"
	"net/http"

	"github.com/adoublef/yap/internal/rooms/sqlite3"
)

func (s *Service) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx = r.Context()
			db  = s.db
		)

		// TODO set session
		rr, err := sqlite3.ListRooms(ctx, db, 10)
		if err != nil {
			log.Println(err)
			// use a http file for errors
			http.Error(w, "Unable to get rooms", http.StatusInternalServerError)
			return
		}

		s.respond(w, r, "index.html", map[string]any{"Rooms": rr})
	}
}
