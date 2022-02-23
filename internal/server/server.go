package server

import (
	"net/http"
	"time"
)

type server struct {
	server *http.Server
}

func NewServer(addr string, handler http.Handler) *server {
	return &server{
		&http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		},
	}
}

func (s *server) Run() error {

	err := s.server.ListenAndServe()

	if err == http.ErrServerClosed {
		return nil
	}

	return err
}
