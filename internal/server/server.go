package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

//TODO заменить логгер на другой
func (s *server) Run() error {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-exit
		s.Stop()

	}()

	fmt.Println("Starting server")
	err := s.server.ListenAndServe()

	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (s *server) Stop() {
	fmt.Println("Stoping server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.server.Shutdown(ctx)

}
