package sqlite3

import (
	"context"
	"database/sql"
	"io/fs"
	"strings"

	"github.com/maragudk/migrate"
	_ "github.com/mattn/go-sqlite3"
)

var (
	args = strings.Join([]string{"_journal=wal", "_timeout=5000", "_synchronous=normal", "_fk=true"}, "&")
)

type DB struct {
	rwc *sql.DB
}

func (db *DB) Conn() *sql.DB {
	if db.rwc == nil {
		panic("rwc cannot be nil")
	}
	return db.rwc
}

func (db *DB) QueryRow(ctx context.Context, query string, args ...any) (row *sql.Row) {
	return db.rwc.QueryRowContext(ctx, query, args...)
}

func (db *DB) Query(ctx context.Context, query string, args ...any) (rows *sql.Rows, err error) {
	return db.rwc.QueryContext(ctx, query, args...)
}

func (db *DB) Begin() (Tx, error) {
	t, err := db.rwc.Begin()
	return &tx{tx: t}, err
}

func (db *DB) Exec(ctx context.Context, query string, args ...any) (result sql.Result, err error) {
	return db.rwc.ExecContext(ctx, query, args...)
}

func Open(ctx context.Context, dsn string) (db *DB, err error) {
	db = &DB{}
	db.rwc, err = sql.Open("sqlite3", dsn+"?"+args)
	if err != nil {
		return nil, err
	}
	// timeout context for ping
	// semi-pointless right now as sqlite will just create a file if non exist
	// ping for a table in the db instead
	// err = db.rwc.PingContext(ctx)
	return
}

// NOTE this is messy, clean-up
type tx struct {
	tx *sql.Tx
}

type Tx interface {
	Rollback() error
	Commit() error
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func (t *tx) Exec(ctx context.Context, query string, args ...any) (result sql.Result, err error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t *tx) Rollback() error {
	return t.tx.Rollback()
}

func (t *tx) Commit() error {
	return t.tx.Commit()
}

type ReadWriter interface {
	Reader
	Writer
}

type Reader interface {
	QueryRow(ctx context.Context, query string, args ...any) (row *sql.Row)
	Query(ctx context.Context, query string, args ...any) (rows *sql.Rows, err error)
}

type Writer interface {
	Begin() (Tx, error)
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type FS struct {
	fsys fs.FS
}

func NewFS(fsys fs.FS) *FS {
	return &FS{fsys: fsys}
}

func (fsys FS) Up(ctx context.Context, db *DB) (err error) {
	return migrate.Up(ctx, db.rwc, fsys.fsys)
}

var (
	ErrNoRows = sql.ErrNoRows
)
