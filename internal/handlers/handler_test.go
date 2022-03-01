package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/siraj18/avito-test/internal/db/postgresdb"
	"github.com/siraj18/avito-test/internal/handlers"
	"github.com/siraj18/avito-test/internal/handlers/mocks"
	"github.com/siraj18/avito-test/internal/models"
	"github.com/siraj18/avito-test/pkg/currencyapi"
	"github.com/stretchr/testify/suite"
)

type handlerSuite struct {
	suite.Suite
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(handlerSuite))
}

func (t *handlerSuite) Test_getBalance() {

	rep := mocks.NewMockRepository()
	rep.On("GetBalance", "f0812ab6-9993-11ec-b909-0242ac120002").Return(&models.User{
		Id:      "f0812ab6-9993-11ec-b909-0242ac120002",
		Balance: 50.0,
	}, nil)

	cur := mocks.NewMockCurrency()
	cur.On("GetCurrency", "USD").Return(25.0, nil) //sad

	h := handlers.NewHandler(rep, cur)

	testSrv := httptest.NewServer(h.InitRoutes())
	client := testSrv.Client()

	body, err := json.Marshal(map[string]string{"id": "f0812ab6-9993-11ec-b909-0242ac120002"})
	t.Nil(err)

	req, err := http.NewRequest("POST", testSrv.URL+"/balance?currency=USD", bytes.NewReader(body))
	t.Nil(err)

	resp, err := client.Do(req)
	t.Nil(err)
	defer resp.Body.Close()

	u := models.User{}
	json.NewDecoder(resp.Body).Decode(&u)

	t.Equal(http.StatusOK, resp.StatusCode)
	t.Equal("f0812ab6-9993-11ec-b909-0242ac120002", u.Id)
	t.Equal(2.0, u.Balance)

}

func (t *handlerSuite) Test_getBalanceUserNotFound() {

	rep := mocks.NewMockRepository()
	rep.On("GetBalance", "f0812ab6-9993-11ec-b909-0242ac120002").Return(nil, postgresdb.ErrorUserNotFound)

	cur := mocks.NewMockCurrency()
	cur.On("GetCurrency", "USD").Return(25.0, nil) //sad

	h := handlers.NewHandler(rep, cur)

	testSrv := httptest.NewServer(h.InitRoutes())
	client := testSrv.Client()

	body, err := json.Marshal(map[string]string{"id": "f0812ab6-9993-11ec-b909-0242ac120002"})
	t.Nil(err)

	req, err := http.NewRequest("POST", testSrv.URL+"/balance?currency=USD", bytes.NewReader(body))
	t.Nil(err)

	resp, err := client.Do(req)
	t.Nil(err)
	defer resp.Body.Close()

	t.Equal(http.StatusNotFound, resp.StatusCode)
}

func (t *handlerSuite) Test_getBalanceSomeError() {

	rep := mocks.NewMockRepository()
	rep.On("GetBalance", "f0812ab6-9993-11ec-b909-0242ac120002").Return(nil, fmt.Errorf("some error"))

	cur := mocks.NewMockCurrency()
	cur.On("GetCurrency", "RUB").Return(25.0, nil) //sad

	h := handlers.NewHandler(rep, cur)

	testSrv := httptest.NewServer(h.InitRoutes())
	client := testSrv.Client()

	body, err := json.Marshal(map[string]string{"id": "f0812ab6-9993-11ec-b909-0242ac120002"})
	t.Nil(err)

	req, err := http.NewRequest("POST", testSrv.URL+"/balance", bytes.NewReader(body))
	t.Nil(err)

	resp, err := client.Do(req)
	t.Nil(err)
	defer resp.Body.Close()

	t.Equal(http.StatusInternalServerError, resp.StatusCode)
}

func (t *handlerSuite) Test_getBalanceWrongCurrency() {
	rep := mocks.NewMockRepository()
	rep.On("GetBalance", "f0812ab6-9993-11ec-b909-0242ac120002").Return(&models.User{
		Id:      "f0812ab6-9993-11ec-b909-0242ac120002",
		Balance: 50.0,
	}, nil)

	cur := mocks.NewMockCurrency()
	cur.On("GetCurrency", "EER").Return(0.0, currencyapi.ErrorWrongCurrency) //sad

	h := handlers.NewHandler(rep, cur)

	testSrv := httptest.NewServer(h.InitRoutes())
	client := testSrv.Client()

	body, err := json.Marshal(map[string]string{"id": "f0812ab6-9993-11ec-b909-0242ac120002"})
	t.Nil(err)

	req, err := http.NewRequest("POST", testSrv.URL+"/balance?currency=EER", bytes.NewReader(body))
	t.Nil(err)

	resp, err := client.Do(req)
	t.Nil(err)

	defer resp.Body.Close()

	t.Equal(http.StatusOK, resp.StatusCode)
}

func (t *handlerSuite) Test_addBalance() {
	rep := mocks.NewMockRepository()
	rep.On("ChangeBalance", "f0812ab6-9993-11ec-b909-0242ac120002", 50.0).Return(&models.User{
		Id:      "f0812ab6-9993-11ec-b909-0242ac120002",
		Balance: 50.0,
	}, nil)

	h := handlers.NewHandler(rep, nil)

	testSrv := httptest.NewServer(h.InitRoutes())
	client := testSrv.Client()

	body, err := json.Marshal(
		map[string]interface{}{
			"id":    "f0812ab6-9993-11ec-b909-0242ac120002",
			"money": 50.0,
		},
	)
	t.Nil(err)

	req, err := http.NewRequest("POST", testSrv.URL+"/changeBalance", bytes.NewReader(body))
	t.Nil(err)

	resp, err := client.Do(req)
	t.Nil(err)
	defer resp.Body.Close()

	u := models.User{}
	json.NewDecoder(resp.Body).Decode(&u)

	t.Equal(u.Id, "f0812ab6-9993-11ec-b909-0242ac120002")
	t.Equal(u.Balance, 50.0)
	t.Equal(http.StatusOK, resp.StatusCode)
}

func (t *handlerSuite) Test_withdrawBalanceNotEnoughMoney() {
	rep := mocks.NewMockRepository()
	rep.On("ChangeBalance", "f0812ab6-9993-11ec-b909-0242ac120002", -1000.0).Return(nil, postgresdb.ErrorNotEnoughMoney)

	h := handlers.NewHandler(rep, nil)

	testSrv := httptest.NewServer(h.InitRoutes())
	client := testSrv.Client()

	body, err := json.Marshal(
		map[string]interface{}{
			"id":    "f0812ab6-9993-11ec-b909-0242ac120002",
			"money": -1000.0,
		},
	)
	t.Nil(err)

	req, err := http.NewRequest("POST", testSrv.URL+"/changeBalance", bytes.NewReader(body))
	t.Nil(err)

	resp, err := client.Do(req)
	t.Nil(err)
	defer resp.Body.Close()

	t.Equal(http.StatusOK, resp.StatusCode)
}

func (t *handlerSuite) Test_addBalanceSomeError() {
	rep := mocks.NewMockRepository()
	rep.On("ChangeBalance", "f0812ab6-9993-11ec-b909-0242ac120002", 50.0).Return(nil, fmt.Errorf("some error"))

	h := handlers.NewHandler(rep, nil)

	testSrv := httptest.NewServer(h.InitRoutes())
	client := testSrv.Client()

	body, err := json.Marshal(
		map[string]interface{}{
			"id":    "f0812ab6-9993-11ec-b909-0242ac120002",
			"money": 50.0,
		},
	)
	t.Nil(err)

	req, err := http.NewRequest("POST", testSrv.URL+"/changeBalance", bytes.NewReader(body))
	t.Nil(err)

	resp, err := client.Do(req)
	t.Nil(err)
	defer resp.Body.Close()

	t.Equal(http.StatusInternalServerError, resp.StatusCode)
}

func (t *handlerSuite) Test_transferBalance() {
	rep := mocks.NewMockRepository()
	rep.On("TransferBalance", "f0812ab6-9993-11ec-b909-0242ac120002", "f0812ab6-9993-11ec-b909-0242ac120003", 50.0).Return(nil)

	h := handlers.NewHandler(rep, nil)

	testSrv := httptest.NewServer(h.InitRoutes())
	client := testSrv.Client()

	body, err := json.Marshal(
		map[string]interface{}{
			"to_id":   "f0812ab6-9993-11ec-b909-0242ac120003",
			"from_id": "f0812ab6-9993-11ec-b909-0242ac120002",
			"money":   50.0,
		},
	)
	t.Nil(err)
	req, err := http.NewRequest("POST", testSrv.URL+"/transferBalance", bytes.NewReader(body))
	t.Nil(err)
	resp, err := client.Do(req)
	t.Nil(err)
	defer resp.Body.Close()

	t.Equal(http.StatusOK, resp.StatusCode)
}

func (t *handlerSuite) Test_transferBalanceNotEnoughMoney() {
	rep := mocks.NewMockRepository()
	rep.On("TransferBalance", "f0812ab6-9993-11ec-b909-0242ac120002", "f0812ab6-9993-11ec-b909-0242ac120003", 50.0).Return(postgresdb.ErrorNotEnoughMoney)

	h := handlers.NewHandler(rep, nil)

	testSrv := httptest.NewServer(h.InitRoutes())
	client := testSrv.Client()

	body, err := json.Marshal(
		map[string]interface{}{
			"to_id":   "f0812ab6-9993-11ec-b909-0242ac120003",
			"from_id": "f0812ab6-9993-11ec-b909-0242ac120002",
			"money":   50.0,
		},
	)
	t.Nil(err)
	req, err := http.NewRequest("POST", testSrv.URL+"/transferBalance", bytes.NewReader(body))
	t.Nil(err)
	resp, err := client.Do(req)
	t.Nil(err)
	defer resp.Body.Close()

	t.Equal(http.StatusOK, resp.StatusCode)
}

func (t *handlerSuite) Test_transferBalanceUserNotFound() {
	rep := mocks.NewMockRepository()
	rep.On("TransferBalance", "f0812ab6-9993-11ec-b909-0242ac120002", "f0812ab6-9993-11ec-b909-0242ac120003", 50.0).Return(postgresdb.ErrorUserNotFound)

	h := handlers.NewHandler(rep, nil)

	testSrv := httptest.NewServer(h.InitRoutes())
	client := testSrv.Client()

	body, err := json.Marshal(
		map[string]interface{}{
			"to_id":   "f0812ab6-9993-11ec-b909-0242ac120003",
			"from_id": "f0812ab6-9993-11ec-b909-0242ac120002",
			"money":   50.0,
		},
	)
	t.Nil(err)
	req, err := http.NewRequest("POST", testSrv.URL+"/transferBalance", bytes.NewReader(body))
	t.Nil(err)
	resp, err := client.Do(req)
	t.Nil(err)
	defer resp.Body.Close()

	t.Equal(http.StatusNotFound, resp.StatusCode)
}

func (t *handlerSuite) Test_transferSomeError() {
	rep := mocks.NewMockRepository()
	rep.On("TransferBalance", "f0812ab6-9993-11ec-b909-0242ac120002", "f0812ab6-9993-11ec-b909-0242ac120003", 50.0).Return(fmt.Errorf("some error"))

	h := handlers.NewHandler(rep, nil)

	testSrv := httptest.NewServer(h.InitRoutes())
	client := testSrv.Client()

	body, err := json.Marshal(
		map[string]interface{}{
			"to_id":   "f0812ab6-9993-11ec-b909-0242ac120003",
			"from_id": "f0812ab6-9993-11ec-b909-0242ac120002",
			"money":   50.0,
		},
	)
	t.Nil(err)
	req, err := http.NewRequest("POST", testSrv.URL+"/transferBalance", bytes.NewReader(body))
	t.Nil(err)
	resp, err := client.Do(req)
	t.Nil(err)
	defer resp.Body.Close()

	t.Equal(http.StatusInternalServerError, resp.StatusCode)
}

func (t *handlerSuite) Test_getTransaction() {

	rep := mocks.NewMockRepository()
	rep.On("GetTransaction", "f0812ab6-9993-11ec-b909-0242ac120002").Return(&models.Transaction{
		Id: "f0812ab6-9993-11ec-b909-0242ac120002",
	}, nil)

	h := handlers.NewHandler(rep, nil)

	testSrv := httptest.NewServer(h.InitRoutes())
	client := testSrv.Client()

	body, err := json.Marshal(map[string]string{"id": "f0812ab6-9993-11ec-b909-0242ac120002"})
	t.Nil(err)

	req, err := http.NewRequest("POST", testSrv.URL+"/transaction", bytes.NewReader(body))
	t.Nil(err)

	resp, err := client.Do(req)
	t.Nil(err)
	defer resp.Body.Close()

	tx := models.Transaction{}
	json.NewDecoder(resp.Body).Decode(&tx)

	t.Equal(http.StatusOK, resp.StatusCode)
	t.Equal("f0812ab6-9993-11ec-b909-0242ac120002", tx.Id)

}
func (t *handlerSuite) Test_getTransactionSomeError() {

	rep := mocks.NewMockRepository()
	rep.On("GetTransaction", "f0812ab6-9993-11ec-b909-0242ac120002").Return(nil, fmt.Errorf("some error"))

	h := handlers.NewHandler(rep, nil)

	testSrv := httptest.NewServer(h.InitRoutes())
	client := testSrv.Client()

	body, err := json.Marshal(map[string]string{"id": "f0812ab6-9993-11ec-b909-0242ac120002"})
	t.Nil(err)

	req, err := http.NewRequest("POST", testSrv.URL+"/transaction", bytes.NewReader(body))
	t.Nil(err)

	resp, err := client.Do(req)
	t.Nil(err)

	defer resp.Body.Close()

	t.Equal(http.StatusInternalServerError, resp.StatusCode)

}
func (t *handlerSuite) Test_getTransactionNotFound() {

	rep := mocks.NewMockRepository()
	rep.On("GetTransaction", "f0812ab6-9993-11ec-b909-0242ac120002").Return(nil, postgresdb.ErrorTransactionNotFound)

	h := handlers.NewHandler(rep, nil)

	testSrv := httptest.NewServer(h.InitRoutes())
	client := testSrv.Client()

	body, err := json.Marshal(map[string]string{"id": "f0812ab6-9993-11ec-b909-0242ac120002"})
	t.Nil(err)

	req, err := http.NewRequest("POST", testSrv.URL+"/transaction", bytes.NewReader(body))
	t.Nil(err)

	resp, err := client.Do(req)
	t.Nil(err)

	defer resp.Body.Close()

	t.Equal(http.StatusNotFound, resp.StatusCode)

}
