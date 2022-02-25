package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type handler struct {
	router     *chi.Mux
	logger     *logrus.Logger
	repository Repository
}

type Repository interface {
	GetBalance(string)
	ChangeBalance(string, float64)
	TransferBalance(string, string, float64)
}

func NewHandler(rep Repository) *handler {
	return &handler{
		router:     chi.NewRouter(),
		logger:     logrus.New(),
		repository: rep,
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
