package mocks

import "github.com/stretchr/testify/mock"

type MockCurrency struct {
	mock.Mock
}

func NewMockCurrency() *MockCurrency {
	return &MockCurrency{}
}

func (m *MockCurrency) GetCurrency(currency string) (float64, error) {
	args := m.Called(currency)
	arg0 := args.Get(0)

	return arg0.(float64), args.Error(1)
}
