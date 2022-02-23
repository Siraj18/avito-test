package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	router *chi.Mux
}

func NewHandler() *handler {
	return &handler{
		router: chi.NewRouter(),
	}
}

// API

func (handler *handler) getBalance(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Balance")
}

func (handler *handler) InitRoutes() *chi.Mux {
	handler.router.Post("/balance", handler.getBalance)
	return handler.router
}
