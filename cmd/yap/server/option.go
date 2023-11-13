package server

import (
	"errors"
	"flag"
	"html/template"
	"os"

	"github.com/go-chi/chi/v5"
)

var (
	dsnRooms = os.Getenv("DATA_SOURCE_ROOMS")
	dsnIam   = os.Getenv("DATA_SOURCE_IAM")
)

type Options struct {
	Addr     string
	RoomsDSN string
	IAMDSN   string
	Router   chi.Router
	Template *template.Template
}

func Parse(fs *flag.FlagSet, args []string) (*Options, error) {
	opts := &Options{
		Router:   chi.NewRouter(),
		Template: template.New(""),
	}

	fs.StringVar(&opts.Addr, "addr", ":8080", "bind listen address")
	// accept env
	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	// TODO check that address is valid
	if opts.IAMDSN = dsnIam; dsnIam == "" {
		return nil, errors.New("dsn required for iam")
	}
	if opts.RoomsDSN = dsnRooms; dsnRooms == "" {
		return nil, errors.New("dsn required for rooms")
	}

	return opts, nil
}
