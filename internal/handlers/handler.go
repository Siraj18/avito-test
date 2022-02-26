package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/siraj18/avito-test/internal/db/postgresdb"
	"github.com/siraj18/avito-test/internal/models"
	"github.com/sirupsen/logrus"
)

type handler struct {
	router     *chi.Mux
	logger     *logrus.Logger
	repository Repository
}

type Repository interface {
	GetBalance(string) (*models.User, error)
	ChangeBalance(string, float64) (*models.User, error)
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
//TODO сдлетаь более нормальную обработку ошибок
//TODO сделать конвертацию валют
func (handler *handler) getBalance(w http.ResponseWriter, r *http.Request) {
	var postData models.UserGetBalanceQuery

	err := json.NewDecoder(r.Body).Decode(&postData)
	if err != nil {
		http.Error(w, "invalid post data", http.StatusBadRequest)
		return
	}

	user, err := handler.repository.GetBalance(postData.Id)

	if err == postgresdb.ErrorUserNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	fmt.Println(err)

	json.NewEncoder(w).Encode(user)
}

func (handler *handler) changeBalance(w http.ResponseWriter, r *http.Request) {
	var postData models.UserChangeBalanceQuery

	err := json.NewDecoder(r.Body).Decode(&postData)
	if err != nil {
		http.Error(w, "invalid post data", http.StatusBadRequest)
		return
	}

	user, err := handler.repository.ChangeBalance(postData.Id, postData.Money)
	if err != nil {
		if err == postgresdb.ErrorNotEnoughMoney {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(500)
	}
	json.NewEncoder(w).Encode(user)
}

func (handler *handler) InitRoutes() *chi.Mux {
	handler.router.Post("/balance", handler.getBalance)
	handler.router.Post("/changeBalance", handler.changeBalance)
	return handler.router
}
