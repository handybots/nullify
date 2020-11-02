package server

import (
	"net/http"
	"unicode/utf8"

	"github.com/handybots/inzerobot/storage"
)

type Server struct {
	move string
	addr string
	stor *storage.DB
}

func New(name, addr string, stor *storage.DB) *Server {
	s := &Server{
		move: "https://t.me/" + name,
		addr: addr,
		stor: stor,
	}

	http.HandleFunc("/", s.index)
	return s
}

func (s *Server) Listen() error {
	return http.ListenAndServe(s.addr, nil)
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = r.URL.Path[1:]
	if utf8.RuneCountInString(r.URL.Path) != 64 {
		http.Redirect(w, r, s.move, http.StatusPermanentRedirect)
		return
	}

	link, err := s.stor.Links.ByString(r.URL.Path)
	if err != nil {
		http.Redirect(w, r, s.move, http.StatusPermanentRedirect)
		return
	}

	go s.stor.Views.Create(storage.View{
		LinkID:    link.ID,
		IP:        r.RemoteAddr,
		UserAgent: r.Header.Get("User-Agent"),
	})

	w.Header().Set("Cache-Control", "no-cache")
	http.Redirect(w, r, link.URL, http.StatusPermanentRedirect)
}
