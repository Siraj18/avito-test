package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type handler struct {
	router *chi.Mux
	logger *logrus.Logger
}

func NewHandler() *handler {
	return &handler{
		router: chi.NewRouter(),
		logger: logrus.New(),
	}
}

// API

func (handler *handler) getBalance(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Balance")
}

//TODO поменять get на post
func (handler *handler) InitRoutes() *chi.Mux {
	handler.router.Get("/balance", handler.getBalance)
	return handler.router
}
