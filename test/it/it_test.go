package it

import (
	"errors"
	"testing"

	"github.com/hetfdex/tiny-bank/internal/domain"
	"github.com/hetfdex/tiny-bank/internal/repository/accountrepo"
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

	svc := service.New(userRepo, accountRepo)

	s.svc = svc
}

func (s *IntegrationTestSuite) TestCreateUser() {
	res, err := s.svc.CreateUser(
		service.CreateUserRequest{
			Name: "joe",
		},
	)

	s.Assert().NotEmpty(res.UserID)
	s.Assert().Nil(err)
}

func (s *IntegrationTestSuite) TestCreateAccount() {
	createUserRes, err := s.svc.CreateUser(
		service.CreateUserRequest{
			Name: "joe",
		},
	)

	s.Assert().Nil(err)

	createAccountRes, err := s.svc.CreateAccount(
		service.CreateAccountRequest(createUserRes),
	)

	s.Assert().NotEmpty(createAccountRes.AccountID)
	s.Assert().Nil(err)
}

func (s *IntegrationTestSuite) TestDeactivate() {
	createUserRes, err := s.svc.CreateUser(
		service.CreateUserRequest{
			Name: "joe",
		},
	)

	s.Assert().Nil(err)

	createAccountRes, err := s.svc.CreateAccount(
		service.CreateAccountRequest(createUserRes),
	)

	s.Assert().Nil(err)

	err = s.svc.DeactivateUser(
		service.DeactivateUserRequest(createUserRes),
	)

	s.Assert().Nil(err)

	balanceRes, err := s.svc.Balance(
		service.BalanceRequest{
			UserID:    createUserRes.UserID,
			AccountID: createAccountRes.AccountID,
		},
	)

	s.Assert().Equal(service.BalanceResponse{}, balanceRes)
	s.Assert().Equal(errors.New("user not active"), err)
}

func (s *IntegrationTestSuite) TestDeposit() {
	createUserRes, err := s.svc.CreateUser(
		service.CreateUserRequest{
			Name: "joe",
		},
	)

	s.Assert().Nil(err)

	createAccountRes, err := s.svc.CreateAccount(
		service.CreateAccountRequest(createUserRes),
	)

	s.Assert().Nil(err)

	depositRes, err := s.svc.Deposit(
		service.DepositRequest{
			UserID:    createUserRes.UserID,
			AccountID: createAccountRes.AccountID,
			Amount:    10,
		},
	)

	s.Assert().Equal(service.DepositResponse{Balance: 10}, depositRes)
	s.Assert().Nil(err)
}

func (s *IntegrationTestSuite) TestWithdraw() {
	createUserRes, err := s.svc.CreateUser(
		service.CreateUserRequest{
			Name: "joe",
		},
	)

	s.Assert().Nil(err)

	createAccountRes, err := s.svc.CreateAccount(
		service.CreateAccountRequest(createUserRes),
	)

	s.Assert().Nil(err)

	depositRes, err := s.svc.Deposit(
		service.DepositRequest{
			UserID:    createUserRes.UserID,
			AccountID: createAccountRes.AccountID,
			Amount:    20,
		},
	)

	s.Assert().Equal(service.DepositResponse{Balance: 20}, depositRes)
	s.Assert().Nil(err)

	withdrawRes, err := s.svc.Withdraw(
		service.WithdrawRequest{
			UserID:    createUserRes.UserID,
			AccountID: createAccountRes.AccountID,
			Amount:    10,
		},
	)

	s.Assert().Equal(service.WithdrawResponse{Balance: 10}, withdrawRes)
	s.Assert().Nil(err)
}

func (s *IntegrationTestSuite) TestTransfer() {
	createJoeUserRes, err := s.svc.CreateUser(
		service.CreateUserRequest{
			Name: "joe",
		},
	)

	s.Assert().Nil(err)

	createJoeAccountRes, err := s.svc.CreateAccount(
		service.CreateAccountRequest(createJoeUserRes),
	)

	s.Assert().Nil(err)

	createMaryUserRes, err := s.svc.CreateUser(
		service.CreateUserRequest{
			Name: "mary",
		},
	)

	s.Assert().Nil(err)

	createMaryAccountRes, err := s.svc.CreateAccount(
		service.CreateAccountRequest(createMaryUserRes),
	)

	s.Assert().Nil(err)

	depositRes, err := s.svc.Deposit(
		service.DepositRequest{
			UserID:    createJoeUserRes.UserID,
			AccountID: createJoeAccountRes.AccountID,
			Amount:    20,
		},
	)

	s.Assert().Equal(service.DepositResponse{Balance: 20}, depositRes)
	s.Assert().Nil(err)

	transferRes, err := s.svc.Transfer(
		service.TransferRequest{
			SenderUserID:      createJoeUserRes.UserID,
			ReceiverUserID:    createMaryUserRes.UserID,
			SenderAccountID:   createJoeAccountRes.AccountID,
			ReceiverAccountID: createMaryAccountRes.AccountID,
			Amount:            10,
		},
	)

	s.Assert().Equal(service.TransferResponse{Balance: 10}, transferRes)
	s.Assert().Nil(err)

	balanceMaryRes, err := s.svc.Balance(
		service.BalanceRequest{
			UserID:    createMaryUserRes.UserID,
			AccountID: createMaryAccountRes.AccountID,
		},
	)

	s.Assert().Equal(service.BalanceResponse{Balance: 10}, balanceMaryRes)
	s.Assert().Nil(err)
}

func (s *IntegrationTestSuite) TestBalance() {
	createUserRes, err := s.svc.CreateUser(
		service.CreateUserRequest{
			Name: "joe",
		},
	)

	s.Assert().Nil(err)

	createAccountRes, err := s.svc.CreateAccount(
		service.CreateAccountRequest(createUserRes),
	)

	s.Assert().Nil(err)

	balanceRes, err := s.svc.Balance(
		service.BalanceRequest{
			UserID:    createUserRes.UserID,
			AccountID: createAccountRes.AccountID,
		},
	)

	s.Assert().Equal(service.BalanceResponse{Balance: 0}, balanceRes)
	s.Assert().Nil(err)
}

func (s *IntegrationTestSuite) TestTransactions() {
	createJoeUserRes, err := s.svc.CreateUser(
		service.CreateUserRequest{
			Name: "joe",
		},
	)

	s.Assert().Nil(err)

	createJoeAccountRes, err := s.svc.CreateAccount(
		service.CreateAccountRequest(createJoeUserRes),
	)

	s.Assert().Nil(err)

	createMaryUserRes, err := s.svc.CreateUser(
		service.CreateUserRequest{
			Name: "mary",
		},
	)

	s.Assert().Nil(err)

	createMaryAccountRes, err := s.svc.CreateAccount(
		service.CreateAccountRequest(createMaryUserRes),
	)

	s.Assert().Nil(err)

	depositRes, err := s.svc.Deposit(
		service.DepositRequest{
			UserID:    createJoeUserRes.UserID,
			AccountID: createJoeAccountRes.AccountID,
			Amount:    20,
		},
	)

	s.Assert().Equal(service.DepositResponse{Balance: 20}, depositRes)
	s.Assert().Nil(err)

	transferRes, err := s.svc.Transfer(
		service.TransferRequest{
			SenderUserID:      createJoeUserRes.UserID,
			ReceiverUserID:    createMaryUserRes.UserID,
			SenderAccountID:   createJoeAccountRes.AccountID,
			ReceiverAccountID: createMaryAccountRes.AccountID,
			Amount:            10,
		},
	)

	s.Assert().Equal(service.TransferResponse{Balance: 10}, transferRes)
	s.Assert().Nil(err)

	balanceMaryRes, err := s.svc.Balance(
		service.BalanceRequest{
			UserID:    createMaryUserRes.UserID,
			AccountID: createMaryAccountRes.AccountID,
		},
	)

	s.Assert().Equal(service.BalanceResponse{Balance: 10}, balanceMaryRes)
	s.Assert().Nil(err)

	transactionsJoeRes, err := s.svc.Transactions(
		service.TransactionsRequest{
			UserID:    createJoeUserRes.UserID,
			AccountID: createJoeAccountRes.AccountID,
		},
	)

	s.Assert().Equal(2, len(transactionsJoeRes.Transactions))

	s.Assert().NotEmpty(transactionsJoeRes.Transactions[0].Timestamp)
	s.Assert().Equal("deposit", transactionsJoeRes.Transactions[0].Operation)
	s.Assert().Equal(20, transactionsJoeRes.Transactions[0].Amount)
	s.Assert().Empty(transactionsJoeRes.Transactions[0].ReceiverUserID)
	s.Assert().Empty(transactionsJoeRes.Transactions[0].SenderUserID)
	s.Assert().Empty(transactionsJoeRes.Transactions[0].ReceiverAccountID)
	s.Assert().Empty(transactionsJoeRes.Transactions[0].SenderAccountID)

	s.Assert().NotEmpty(transactionsJoeRes.Transactions[1].Timestamp)
	s.Assert().Equal("transfer", transactionsJoeRes.Transactions[1].Operation)
	s.Assert().Equal(10, transactionsJoeRes.Transactions[1].Amount)
	s.Assert().Equal(createMaryUserRes.UserID, transactionsJoeRes.Transactions[1].ReceiverUserID)
	s.Assert().Empty(transactionsJoeRes.Transactions[1].SenderUserID)
	s.Assert().Equal(createMaryAccountRes.AccountID, transactionsJoeRes.Transactions[1].ReceiverAccountID)
	s.Assert().Empty(transactionsJoeRes.Transactions[1].SenderAccountID)

	s.Assert().Nil(err)

	transactionsMaryRes, err := s.svc.Transactions(
		service.TransactionsRequest{
			UserID:    createMaryUserRes.UserID,
			AccountID: createMaryAccountRes.AccountID,
		},
	)

	s.Assert().Equal(1, len(transactionsMaryRes.Transactions))

	s.Assert().NotEmpty(transactionsMaryRes.Transactions[0].Timestamp)
	s.Assert().Equal("transfer", transactionsMaryRes.Transactions[0].Operation)
	s.Assert().Equal(10, transactionsMaryRes.Transactions[0].Amount)
	s.Assert().Empty(transactionsMaryRes.Transactions[0].ReceiverUserID)
	s.Assert().Equal(createJoeUserRes.UserID, transactionsMaryRes.Transactions[0].SenderUserID)
	s.Assert().Empty(transactionsMaryRes.Transactions[0].ReceiverAccountID)
	s.Assert().Equal(createJoeAccountRes.AccountID, transactionsMaryRes.Transactions[0].SenderAccountID)

	s.Assert().Nil(err)
}
