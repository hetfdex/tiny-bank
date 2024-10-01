package servicemock

import (
	"github.com/hetfdex/tiny-bank/internal/service"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) Create(req service.CreateRequest) (service.CreateResponse, error) {
	args := m.Called(req)

	return args.Get(0).(service.CreateResponse), args.Error(1)
}

func (m *Mock) Deactivate(req service.DeactivateRequest) error {
	args := m.Called(req)

	return args.Error(0)
}

func (m *Mock) Deposit(req service.DepositRequest) (service.DepositResponse, error) {
	args := m.Called(req)

	return args.Get(0).(service.DepositResponse), args.Error(1)
}

func (m *Mock) Withdraw(req service.WithdrawRequest) (service.WithdrawResponse, error) {
	args := m.Called(req)

	return args.Get(0).(service.WithdrawResponse), args.Error(1)
}

func (m *Mock) Transfer(req service.TransferRequest) (service.TransferResponse, error) {
	args := m.Called(req)

	return args.Get(0).(service.TransferResponse), args.Error(1)
}

func (m *Mock) Balance(req service.BalanceRequest) (service.BalanceResponse, error) {
	args := m.Called(req)

	return args.Get(0).(service.BalanceResponse), args.Error(1)
}

func (m *Mock) History(req service.HistoryRequest) (service.HistoryResponse, error) {
	args := m.Called(req)

	return args.Get(0).(service.HistoryResponse), args.Error(1)
}
