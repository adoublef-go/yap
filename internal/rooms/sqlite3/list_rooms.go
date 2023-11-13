package sqlite3

import (
	"context"

	"github.com/adoublef/yap/internal/rooms"
	"github.com/adoublef/yap/sqlite3"
)

func ListRooms(ctx context.Context, db sqlite3.Reader, limit int) ([]*rooms.Room, error) {
	var (
		qry = `
		SELECT r.id
		FROM rooms r
		LIMIT ?`
	)
	r, err := db.Query(ctx, qry, limit)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var rr []*rooms.Room
	for r.Next() {
		var room rooms.Room
		err = r.Scan(&room.ID)
		if err != nil {
			return nil, err
		}
		rr = append(rr, &room)
	}

	if err := r.Err(); err != nil {
		return nil, err
	}
	return rr, nil
}
