package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/siraj18/avito-test/internal/currency"
	"github.com/siraj18/avito-test/internal/db/postgresdb"
	"github.com/siraj18/avito-test/internal/models"
	"github.com/sirupsen/logrus"
)

type handler struct {
	router     *chi.Mux
	logger     *logrus.Logger
	repository Repository
	currency   *currency.Currency
}

type Repository interface {
	GetBalance(string) (*models.User, error)
	ChangeBalance(string, float64) (*models.User, error)
	TransferBalance(string, string, float64) error
	GetTransaction(string) (*models.Transaction, error)
	GetAllTransactions(string, string, int, int) (*[]models.Transaction, error)
}

func NewHandler(rep Repository, currency *currency.Currency) *handler {
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
	defer r.Body.Close()

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
	defer r.Body.Close()

	if err != nil {
		http.Error(w, "invalid post data", http.StatusBadRequest)
		return
	}

	user, err := handler.repository.ChangeBalance(postData.Id, postData.Money)
	if err != nil {
		if err == postgresdb.ErrorNotEnoughMoney {
			http.Error(w, err.Error(), http.StatusOK)
			return
		}
		w.WriteHeader(500)
		handler.logger.Error(err)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (handler *handler) transferBalance(w http.ResponseWriter, r *http.Request) {
	var postData models.UserTransferBalanceQuery

	err := json.NewDecoder(r.Body).Decode(&postData)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, "invalid post data", http.StatusBadRequest)
		return
	}

	err = handler.repository.TransferBalance(postData.FromId, postData.ToId, postData.Money)

	if err != nil {
		switch err {
		case postgresdb.ErrorNotEnoughMoney:
			http.Error(w, err.Error(), http.StatusOK)
		case postgresdb.ErrorUserNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			w.WriteHeader(500)
		}
		handler.logger.Error(err)
		return
	}
	w.WriteHeader(200)
	fmt.Fprintf(w, "The transfer was completed successfully")
}

func (handler *handler) getTransaction(w http.ResponseWriter, r *http.Request) {
	var postData models.TransactionGetQuery

	err := json.NewDecoder(r.Body).Decode(&postData)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, "invalid post data", http.StatusBadRequest)
		return
	}

	transaction, err := handler.repository.GetTransaction(postData.Id)
	if err != nil {
		if err == postgresdb.ErrorTransactionNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(500)
		handler.logger.Error(err)
		return
	}

	json.NewEncoder(w).Encode(transaction)

}

func (handler *handler) getAllTransactions(w http.ResponseWriter, r *http.Request) {
	var postData models.AllTransactionsGetQuery

	err := json.NewDecoder(r.Body).Decode(&postData)
	defer r.Body.Close()

	if err != nil {
		handler.logger.Error(err)
		http.Error(w, "invalid post data", http.StatusBadRequest)
		return
	}

	transactions, err := handler.repository.GetAllTransactions(postData.Id, postData.SortType, postData.Limit, postData.Page)
	if err != nil {
		handler.logger.Error(err)
		w.WriteHeader(500)
		return
	}

	json.NewEncoder(w).Encode(transactions)

}

func (handler *handler) InitRoutes() *chi.Mux {
	handler.router.Post("/balance", handler.getBalance)
	handler.router.Post("/changeBalance", handler.changeBalance)
	handler.router.Post("/transferBalance", handler.transferBalance)
	handler.router.Post("/transaction", handler.getTransaction)
	handler.router.Post("/allTransactions", handler.getAllTransactions)
	return handler.router
}
