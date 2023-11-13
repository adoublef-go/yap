package http

import (
	"embed"
	"html/template"
	"net/http"

	sql "github.com/adoublef/yap/sqlite3"
	"github.com/go-chi/chi/v5"
)

var (
	//go:embed all:*.html
	fsys embed.FS
)

func ParseFS(t *template.Template) (*template.Template, error) {
	return t.ParseFS(fsys, "*.html")
}

type Service struct {
	m  *chi.Mux
	db sql.ReadWriter
	t  *template.Template
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.m.ServeHTTP(w, r)
}

func NewService(db sql.ReadWriter, t *template.Template) (s *Service) {
	s = &Service{
		m:  chi.NewMux(),
		db: db,
		t:  t,
	}
	s.routes()
	return
}

func (s *Service) routes() {
	s.m.Get("/", s.handleIndex())
}

func (s *Service) respond(w http.ResponseWriter, r *http.Request, name string, data any) {
	if err := s.t.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "Failed to render view", http.StatusInternalServerError)
		return
	}
}
