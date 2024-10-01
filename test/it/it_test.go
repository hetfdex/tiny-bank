package it

import (
	"errors"
	"testing"

	"github.com/hetfdex/tiny-bank/internal/domain"
	"github.com/hetfdex/tiny-bank/internal/repository/accountrepo"
	"github.com/hetfdex/tiny-bank/internal/repository/historyrepo"
	"github.com/hetfdex/tiny-bank/internal/repository/userrepo"
	"github.com/hetfdex/tiny-bank/internal/service"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	svc service.Service
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite() {
	userRepo := userrepo.New(make(map[string]domain.User))
	accountRepo := accountrepo.New(make(map[string]domain.Account))
	historyRepo := historyrepo.New(make(map[string]domain.History))

	svc := service.New(userRepo, accountRepo, historyRepo)

	s.svc = svc
}

func (s *IntegrationTestSuite) TestCreate() {
	res, err := s.svc.Create(
		service.CreateRequest{
			Name: "joe",
		},
	)

	s.Assert().NotEmpty(res.ID)
	s.Assert().Equal(1, len(res.AccountIDs))

	s.Assert().Nil(err)
}

func (s *IntegrationTestSuite) TestDeactivate() {
	createRes, err := s.svc.Create(
		service.CreateRequest{
			Name: "joe",
		},
	)

	s.Assert().Nil(err)

	err = s.svc.Deactivate(
		service.DeactivateRequest{
			UserID: createRes.ID,
		},
	)

	s.Assert().Nil(err)

	balanceRes, err := s.svc.Balance(
		service.BalanceRequest{
			UserID:    createRes.ID,
			AccountID: createRes.AccountIDs[0],
		},
	)

	s.Assert().Equal(service.BalanceResponse{}, balanceRes)
	s.Assert().Equal(errors.New("user not active"), err)
}

func (s *IntegrationTestSuite) TestDeposit() {
	createRes, err := s.svc.Create(
		service.CreateRequest{
			Name: "joe",
		},
	)

	s.Assert().Nil(err)

	depositRes, err := s.svc.Deposit(
		service.DepositRequest{
			UserID:    createRes.ID,
			AccountID: createRes.AccountIDs[0],
			Amount:    10,
		},
	)

	s.Assert().Equal(service.DepositResponse{Balance: 10}, depositRes)
	s.Assert().Nil(err)
}

func (s *IntegrationTestSuite) TestWithdraw() {
	createRes, err := s.svc.Create(
		service.CreateRequest{
			Name: "joe",
		},
	)

	s.Assert().Nil(err)

	depositRes, err := s.svc.Deposit(
		service.DepositRequest{
			UserID:    createRes.ID,
			AccountID: createRes.AccountIDs[0],
			Amount:    20,
		},
	)

	s.Assert().Equal(service.DepositResponse{Balance: 20}, depositRes)
	s.Assert().Nil(err)

	withdrawRes, err := s.svc.Withdraw(
		service.WithdrawRequest{
			UserID:    createRes.ID,
			AccountID: createRes.AccountIDs[0],
			Amount:    10,
		},
	)

	s.Assert().Equal(service.WithdrawResponse{Balance: 10}, withdrawRes)
	s.Assert().Nil(err)
}

func (s *IntegrationTestSuite) TestTransfer() {
	createJoeRes, err := s.svc.Create(
		service.CreateRequest{
			Name: "joe",
		},
	)

	s.Assert().Nil(err)

	createMaryRes, err := s.svc.Create(
		service.CreateRequest{
			Name: "mary",
		},
	)

	s.Assert().Nil(err)

	depositRes, err := s.svc.Deposit(
		service.DepositRequest{
			UserID:    createJoeRes.ID,
			AccountID: createJoeRes.AccountIDs[0],
			Amount:    20,
		},
	)

	s.Assert().Equal(service.DepositResponse{Balance: 20}, depositRes)
	s.Assert().Nil(err)

	transferRes, err := s.svc.Transfer(
		service.TransferRequest{
			SenderUserID:      createJoeRes.ID,
			ReceiverUserID:    createMaryRes.ID,
			SenderAccountID:   createJoeRes.AccountIDs[0],
			ReceiverAccountID: createMaryRes.AccountIDs[0],
			Amount:            10,
		},
	)

	s.Assert().Equal(service.TransferResponse{Balance: 10}, transferRes)
	s.Assert().Nil(err)

	balanceMaryRes, err := s.svc.Balance(
		service.BalanceRequest{
			UserID:    createMaryRes.ID,
			AccountID: createMaryRes.AccountIDs[0],
		},
	)

	s.Assert().Equal(service.BalanceResponse{Balance: 10}, balanceMaryRes)
	s.Assert().Nil(err)
}

func (s *IntegrationTestSuite) TestBalance() {
	createRes, err := s.svc.Create(
		service.CreateRequest{
			Name: "joe",
		},
	)

	s.Assert().Nil(err)

	balanceRes, err := s.svc.Balance(
		service.BalanceRequest{
			UserID:    createRes.ID,
			AccountID: createRes.AccountIDs[0],
		},
	)

	s.Assert().Equal(service.BalanceResponse{Balance: 0}, balanceRes)
	s.Assert().Nil(err)
}

func (s *IntegrationTestSuite) TestHistory() {
	createJoeRes, err := s.svc.Create(
		service.CreateRequest{
			Name: "joe",
		},
	)

	s.Assert().Nil(err)

	createMaryRes, err := s.svc.Create(
		service.CreateRequest{
			Name: "mary",
		},
	)

	s.Assert().Nil(err)

	depositRes, err := s.svc.Deposit(
		service.DepositRequest{
			UserID:    createJoeRes.ID,
			AccountID: createJoeRes.AccountIDs[0],
			Amount:    20,
		},
	)

	s.Assert().Equal(service.DepositResponse{Balance: 20}, depositRes)
	s.Assert().Nil(err)

	transferRes, err := s.svc.Transfer(
		service.TransferRequest{
			SenderUserID:      createJoeRes.ID,
			ReceiverUserID:    createMaryRes.ID,
			SenderAccountID:   createJoeRes.AccountIDs[0],
			ReceiverAccountID: createMaryRes.AccountIDs[0],
			Amount:            10,
		},
	)

	s.Assert().Equal(service.TransferResponse{Balance: 10}, transferRes)
	s.Assert().Nil(err)

	balanceMaryRes, err := s.svc.Balance(
		service.BalanceRequest{
			UserID:    createMaryRes.ID,
			AccountID: createMaryRes.AccountIDs[0],
		},
	)

	s.Assert().Equal(service.BalanceResponse{Balance: 10}, balanceMaryRes)
	s.Assert().Nil(err)

	historyJoeRes, err := s.svc.History(
		service.HistoryRequest{
			UserID:    createJoeRes.ID,
			AccountID: createJoeRes.AccountIDs[0],
		},
	)

	s.Assert().Equal(2, len(historyJoeRes.Events))

	s.Assert().NotEmpty(historyJoeRes.Events[0].Timestamp)
	s.Assert().Equal("deposit", historyJoeRes.Events[0].Operation)
	s.Assert().Equal(20, historyJoeRes.Events[0].Amount)
	s.Assert().Empty(historyJoeRes.Events[0].ReceiverUserID)
	s.Assert().Empty(historyJoeRes.Events[0].SenderUserID)
	s.Assert().Empty(historyJoeRes.Events[0].ReceiverAccountID)
	s.Assert().Empty(historyJoeRes.Events[0].SenderAccountID)

	s.Assert().NotEmpty(historyJoeRes.Events[1].Timestamp)
	s.Assert().Equal("transfer", historyJoeRes.Events[1].Operation)
	s.Assert().Equal(10, historyJoeRes.Events[1].Amount)
	s.Assert().Equal(createMaryRes.ID, historyJoeRes.Events[1].ReceiverUserID)
	s.Assert().Empty(historyJoeRes.Events[1].SenderUserID)
	s.Assert().Equal(createMaryRes.AccountIDs[0], historyJoeRes.Events[1].ReceiverAccountID)
	s.Assert().Empty(historyJoeRes.Events[1].SenderAccountID)

	s.Assert().Nil(err)

	historyMaryRes, err := s.svc.History(
		service.HistoryRequest{
			UserID:    createMaryRes.ID,
			AccountID: createMaryRes.AccountIDs[0],
		},
	)

	s.Assert().Equal(1, len(historyMaryRes.Events))

	s.Assert().NotEmpty(historyMaryRes.Events[0].Timestamp)
	s.Assert().Equal("transfer", historyMaryRes.Events[0].Operation)
	s.Assert().Equal(10, historyMaryRes.Events[0].Amount)
	s.Assert().Empty(historyMaryRes.Events[0].ReceiverUserID)
	s.Assert().Equal(createJoeRes.ID, historyMaryRes.Events[0].SenderUserID)
	s.Assert().Empty(historyMaryRes.Events[0].ReceiverAccountID)
	s.Assert().Equal(createJoeRes.AccountIDs[0], historyMaryRes.Events[0].SenderAccountID)

	s.Assert().Nil(err)
}
