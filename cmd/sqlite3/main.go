package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/adoublef/yap/cmd/yap/server"
	iam "github.com/adoublef/yap/internal/iam/sqlite3"
	rooms "github.com/adoublef/yap/internal/rooms/sqlite3"
	sql "github.com/adoublef/yap/sqlite3"
)

var (
	timeout = 5 * time.Second
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// prepare server options
	fs := flag.NewFlagSet("", flag.ExitOnError)
	// TODO set usage
	opts, err := server.Parse(fs, os.Args[1:])
	if err != nil {
		log.Fatalln(err)
	}
	if err := run(ctx, opts); err != nil {
		log.Fatalln(err)
	}
}

func run(ctx context.Context, opts *server.Options) (err error) {
	// prepare rooms
	{
		db, err := sql.Open(ctx, opts.RoomsDSN)
		if err != nil {
			return err
		}

		if err = rooms.Migrate.Up(ctx, db); err != nil {
			return err
		}
	}
	// prepare iam
	{
		db, err := sql.Open(ctx, opts.IAMDSN)
		if err != nil {
			return err
		}

		if err = iam.Migrate.Up(ctx, db); err != nil {
			return err
		}
	}
	return nil
}
