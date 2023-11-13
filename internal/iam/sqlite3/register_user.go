package sqlite3

import (
	"context"

	"github.com/adoublef/yap/internal/iam"
	"github.com/adoublef/yap/sqlite3"
)

func RegisterUser(ctx context.Context, db sqlite3.Writer, u *iam.User) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(ctx, `INSERT INTO profiles (id, login, name, photo)
	VALUES (?, ?, ?, ?)`, u.Profile.ID, u.Profile.Login, u.Profile.Name, u.Profile.Photo)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `INSERT INTO accounts (id, profile, email)
	VALUES (?, ?, ?)`, u.OAuth2, u.Profile.ID, u.Email)
	if err != nil {
		return err
	}
	return tx.Commit()
}
