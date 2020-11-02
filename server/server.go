package server

import (
	"net/http"
	"unicode/utf8"

	"github.com/handybots/inzerobot/storage"
)

type Server struct {
	addr string
	stor *storage.DB
}

func NewServer(address string, stor *storage.DB) *Server {
	s := &Server{
		addr: address,
		stor: stor,
	}

	http.HandleFunc("/", s.index)
	return s
}

func (s *Server) Listen() error {
	return http.ListenAndServe(s.addr, nil)
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	if utf8.RuneCountInString(r.URL.Path) != 64 {
		r.URL.Path = ""
		http.Redirect(w, r, r.URL.String(), http.StatusPermanentRedirect)
		return
	}

	link, err := s.stor.Links.ByString(r.URL.Path)
	if err != nil {
		r.URL.Path = ""
		http.Redirect(w, r, r.URL.String(), http.StatusPermanentRedirect)
		return
	}

	s.stor.Views.Create(storage.View{
		LinkID:    link.ID,
		IP:        r.RemoteAddr,
		UserAgent: r.Header.Get("User-Agent"),
	})

	http.Redirect(w, r, link.URL, http.StatusPermanentRedirect)
}
