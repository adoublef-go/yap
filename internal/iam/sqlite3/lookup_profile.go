package sqlite3

// https://stackoverflow.com/questions/418898/upsert-not-insert-or-replace
import (
	"context"

	"github.com/adoublef/yap/internal/iam"
	sql "github.com/adoublef/yap/sqlite3"
)

func LookupProfile(ctx context.Context, db sql.Reader, acc string) (profile *iam.Profile, err error) {
	var (
		qry = `
		SELECT p.id, p.login, p.name, p.photo
		FROM profiles p
		INNER JOIN accounts a 
			ON a.profile = p.id 
		WHERE a.id = ?`
	)
	profile = &iam.Profile{}
	err = db.QueryRow(ctx, qry, acc).Scan(
		&profile.ID,
		&profile.Login,
		&profile.Name,
		&profile.Photo,
	)
	return
}
