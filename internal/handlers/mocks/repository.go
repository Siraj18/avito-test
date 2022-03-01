package mocks

import (
	"github.com/siraj18/avito-test/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func NewMockRepository() *MockRepository {
	return &MockRepository{}
}

func (m *MockRepository) GetBalance(id string) (*models.User, error) {
	args := m.Called(id)

	arg0 := args.Get(0)
	if arg0 == nil {
		return nil, args.Error(1)
	}

	return arg0.(*models.User), args.Error(1)
}

func (m *MockRepository) ChangeBalance(id string, money float64) (*models.User, error) {
	args := m.Called(id, money)

	arg0 := args.Get(0)
	if arg0 == nil {
		return nil, args.Error(1)
	}

	return arg0.(*models.User), args.Error(1)

}

func (m *MockRepository) TransferBalance(fromId, toId string, money float64) error {
	args := m.Called(fromId, toId, money)

	return args.Error(0)
}

func (m *MockRepository) GetTransaction(id string) (*models.Transaction, error) {
	args := m.Called(id)

	arg0 := args.Get(0)
	if arg0 == nil {
		return nil, args.Error(1)
	}

	return arg0.(*models.Transaction), args.Error(1)
}

func (m *MockRepository) GetAllTransactions(id, sortType string, limit, page int) (*[]models.Transaction, error) {
	args := m.Called(id, sortType, limit, page)

	arg0 := args.Get(0)

	if arg0 == nil {
		return nil, args.Error(1)
	}

	return arg0.(*[]models.Transaction), args.Error(1)
}
