package server

import (
	"context"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"time"

	httputil "github.com/adoublef/yap/http"
	iam "github.com/adoublef/yap/internal/iam/http"
	rooms "github.com/adoublef/yap/internal/rooms/http"
	natsutil "github.com/adoublef/yap/nats"
	sql "github.com/adoublef/yap/sqlite3"
	"github.com/adoublef/yap/static"
	"github.com/go-chi/chi/v5"
)

var (
	timeout = 5 * time.Second
)

type Server struct {
	s *http.Server
}

func NewServer(ctx context.Context, nc *natsutil.Conn, opts *Options) (s *Server, err error) {
	mux, t := chi.NewMux(), opts.Template.Funcs(funcMap).Funcs(static.FuncMap)
	// prepare http server
	mux.Handle("/static/*", http.StripPrefix("/static/", static.Handler))
	users, err := httputil.NewCookieJar(nc, "user-sessions", 10*time.Minute)
	if err != nil {
		return nil, err
	}
	// prepare iam
	{
		// prepare sqlite
		db, err := sql.Open(ctx, opts.IAMDSN)
		if err != nil {
			return nil, fmt.Errorf("open database: %w", err)
		}
		mux.Mount("/", iam.NewService(db, users))
	}
	// prepare rooms
	{
		// prepare sqlite
		db, err := sql.Open(ctx, opts.RoomsDSN)
		if err != nil {
			return nil, fmt.Errorf("open database: %w", err)
		}
		// prepare template
		if t, err = rooms.ParseFS(t); err != nil {
			return nil, fmt.Errorf("parse template: %w", err)
		}
		mux.Mount("/rooms", rooms.NewService(db, t))
	}
	s = &Server{
		&http.Server{
			Addr:    opts.Addr,
			Handler: mux,
			BaseContext: func(l net.Listener) context.Context {
				return ctx
			},
		},
	}
	return
}

func (s *Server) ListenAndServe() (err error) {
	err = s.s.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) (err error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return s.s.Shutdown(ctx)
}

var funcMap = template.FuncMap{
	"env": func(key string) string {
		return os.Getenv(key)
	},
}
