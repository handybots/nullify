package server

import (
	"fmt"
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
	if r.URL.Path == "/" {
		fmt.Fprint(w, "hello")
		return
	}

	r.URL.Path = r.URL.Path[1:]
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

	go s.stor.Views.Create(storage.View{
		LinkID:    link.ID,
		IP:        r.RemoteAddr,
		UserAgent: r.Header.Get("User-Agent"),
	})

	w.Header().Set("Cache-Control", "no-cache")
	http.Redirect(w, r, link.URL, http.StatusPermanentRedirect)
}
