package service

import (
	"errors"
	"time"

	guuid "github.com/google/uuid"
	"github.com/hetfdex/tiny-bank/internal/domain"
	"github.com/hetfdex/tiny-bank/internal/repository/accountrepo"
	"github.com/hetfdex/tiny-bank/internal/repository/historyrepo"
	"github.com/hetfdex/tiny-bank/internal/repository/userrepo"
)

type Service interface {
	Create(CreateRequest) (CreateResponse, error)
	Deactivate(DeactivateRequest) error
	Deposit(DepositRequest) (DepositResponse, error)
	Withdraw(WithdrawRequest) (WithdrawResponse, error)
	Transfer(TransferRequest) (TransferResponse, error)
	Balance(BalanceRequest) (BalanceResponse, error)
	History(HistoryRequest) (HistoryResponse, error)
}

type svc struct {
	userRepo    userrepo.Repo
	accountRepo accountrepo.Repo
	historyRepo historyrepo.Repo
}

func New(
	userRepo userrepo.Repo,
	accountRepo accountrepo.Repo,
	historyRepo historyrepo.Repo,
) Service {
	return &svc{
		userRepo:    userRepo,
		accountRepo: accountRepo,
		historyRepo: historyRepo,
	}
}

func (s svc) Create(req CreateRequest) (CreateResponse, error) {
	if req.Name == "" {
		return CreateResponse{}, errors.New("invalid user name")
	}

	history, err := s.historyRepo.Create(historyrepo.CreateRequest{})

	if err != nil {
		return CreateResponse{}, err
	}

	account, err := s.accountRepo.Create(
		accountrepo.CreateRequest{
			HistoryID: history.ID,
		},
	)

	if err != nil {
		return CreateResponse{}, err
	}

	user, err := s.userRepo.Create(
		userrepo.CreateRequest{
			Name:      req.Name,
			AccountID: account.ID,
		},
	)

	if err != nil {
		return CreateResponse{}, err
	}

	return CreateResponse{
		ID:         user.ID,
		AccountIDs: user.AccountIDs,
	}, nil
}

func (s svc) Deactivate(req DeactivateRequest) error {
	if !validID(req.UserID) {
		return errors.New("invalid user id")
	}

	return s.userRepo.Update(
		userrepo.UpdateRequest{
			ID:     req.UserID,
			Active: false,
		},
	)
}

func (s svc) Deposit(req DepositRequest) (DepositResponse, error) {
	if !validID(req.UserID) {
		return DepositResponse{}, errors.New("invalid user id")
	}

	if !validID(req.AccountID) {
		return DepositResponse{}, errors.New("invalid account id")
	}

	if req.Amount <= 0 {
		return DepositResponse{}, errors.New("invalid amount")
	}

	user, err := s.userRepo.Read(
		userrepo.ReadRequest{
			ID: req.UserID,
		},
	)

	if err != nil {
		return DepositResponse{}, err
	}

	if !userAccount(user.AccountIDs, req.AccountID) {
		return DepositResponse{}, errors.New("unauthorized account id")
	}

	account, err := s.accountRepo.Read(
		accountrepo.ReadRequest{
			ID: req.AccountID,
		},
	)

	if err != nil {
		return DepositResponse{}, err
	}
	balance := account.Balance + req.Amount

	err = s.accountRepo.Update(
		accountrepo.UpdateRequest{
			ID:      req.AccountID,
			Balance: balance,
		},
	)

	if err != nil {
		return DepositResponse{}, err
	}

	err = s.historyRepo.Update(
		historyrepo.UpdateRequest{
			ID: account.HistoryID,
			Event: domain.Event{
				Timestamp: time.Now().UTC(),
				Operation: "deposit",
				Amount:    req.Amount,
			},
		},
	)

	if err != nil {
		return DepositResponse{}, err
	}

	return DepositResponse{
		Balance: balance,
	}, nil
}

func (s svc) Withdraw(req WithdrawRequest) (WithdrawResponse, error) {
	if !validID(req.UserID) {
		return WithdrawResponse{}, errors.New("invalid user id")
	}

	if !validID(req.AccountID) {
		return WithdrawResponse{}, errors.New("invalid account id")
	}

	if req.Amount <= 0 {
		return WithdrawResponse{}, errors.New("invalid amount")
	}

	user, err := s.userRepo.Read(
		userrepo.ReadRequest{
			ID: req.UserID,
		},
	)

	if err != nil {
		return WithdrawResponse{}, err
	}

	if !userAccount(user.AccountIDs, req.AccountID) {
		return WithdrawResponse{}, errors.New("unauthorized account id")
	}

	account, err := s.accountRepo.Read(
		accountrepo.ReadRequest{
			ID: req.AccountID,
		},
	)

	if err != nil {
		return WithdrawResponse{}, err
	}

	if account.Balance < req.Amount {
		return WithdrawResponse{}, errors.New("insuficient funds")
	}

	balance := account.Balance - req.Amount

	err = s.accountRepo.Update(
		accountrepo.UpdateRequest{
			ID:      req.AccountID,
			Balance: balance,
		},
	)

	if err != nil {
		return WithdrawResponse{}, err
	}

	err = s.historyRepo.Update(
		historyrepo.UpdateRequest{
			ID: account.HistoryID,
			Event: domain.Event{
				Timestamp: time.Now().UTC(),
				Operation: "withdraw",
				Amount:    req.Amount,
			},
		},
	)

	if err != nil {
		return WithdrawResponse{}, err
	}

	return WithdrawResponse{
		Balance: balance,
	}, nil
}

func (s svc) Transfer(req TransferRequest) (TransferResponse, error) {
	if !validID(req.SenderUserID) {
		return TransferResponse{}, errors.New("invalid sender user id")
	}

	if !validID(req.ReceiverUserID) {
		return TransferResponse{}, errors.New("invalid receiver user id")
	}

	if !validID(req.SenderAccountID) {
		return TransferResponse{}, errors.New("invalid sender account id")
	}

	if !validID(req.ReceiverAccountID) {
		return TransferResponse{}, errors.New("invalid receiver account id")
	}

	if req.Amount <= 0 {
		return TransferResponse{}, errors.New("invalid amount")
	}

	if req.SenderAccountID == req.ReceiverAccountID {
		return TransferResponse{}, errors.New("same account")
	}

	sender, err := s.userRepo.Read(
		userrepo.ReadRequest{
			ID: req.SenderUserID,
		},
	)

	if err != nil {
		return TransferResponse{}, err
	}

	if !userAccount(sender.AccountIDs, req.SenderAccountID) {
		return TransferResponse{}, errors.New("unauthorized account id")
	}

	senderAccount, err := s.accountRepo.Read(
		accountrepo.ReadRequest{
			ID: req.SenderAccountID,
		},
	)

	if err != nil {
		return TransferResponse{}, err
	}

	if senderAccount.Balance < req.Amount {
		return TransferResponse{}, errors.New("insuficient funds")
	}

	receiver, err := s.userRepo.Read(
		userrepo.ReadRequest{
			ID: req.ReceiverUserID,
		},
	)

	if err != nil {
		return TransferResponse{}, err
	}

	if !userAccount(receiver.AccountIDs, req.ReceiverAccountID) {
		return TransferResponse{}, errors.New("unauthorized account id")
	}

	receiverAccount, err := s.accountRepo.Read(
		accountrepo.ReadRequest{
			ID: req.ReceiverAccountID,
		},
	)

	if err != nil {
		return TransferResponse{}, err
	}

	senderBalance := senderAccount.Balance - req.Amount

	err = s.accountRepo.Update(
		accountrepo.UpdateRequest{
			ID:      req.SenderAccountID,
			Balance: senderBalance,
		},
	)

	if err != nil {
		return TransferResponse{}, err
	}

	err = s.accountRepo.Update(
		accountrepo.UpdateRequest{
			ID:      req.ReceiverAccountID,
			Balance: receiverAccount.Balance + req.Amount,
		},
	)

	if err != nil {
		return TransferResponse{}, err
	}

	err = s.historyRepo.Update(
		historyrepo.UpdateRequest{
			ID: senderAccount.HistoryID,
			Event: domain.Event{
				Timestamp:         time.Now().UTC(),
				Operation:         "transfer",
				Amount:            req.Amount,
				ReceiverUserID:    req.ReceiverUserID,
				ReceiverAccountID: req.ReceiverAccountID,
			},
		},
	)

	if err != nil {
		return TransferResponse{}, err
	}

	err = s.historyRepo.Update(
		historyrepo.UpdateRequest{
			ID: receiverAccount.HistoryID,
			Event: domain.Event{
				Timestamp:       time.Now().UTC(),
				Operation:       "transfer",
				Amount:          req.Amount,
				SenderUserID:    req.SenderUserID,
				SenderAccountID: req.SenderAccountID,
			},
		},
	)

	if err != nil {
		return TransferResponse{}, err
	}

	return TransferResponse{
		Balance: senderBalance,
	}, nil
}

func (s svc) Balance(req BalanceRequest) (BalanceResponse, error) {
	if !validID(req.UserID) {
		return BalanceResponse{}, errors.New("invalid user id")
	}

	if !validID(req.AccountID) {
		return BalanceResponse{}, errors.New("invalid account id")
	}

	user, err := s.userRepo.Read(
		userrepo.ReadRequest{
			ID: req.UserID,
		},
	)

	if err != nil {
		return BalanceResponse{}, err
	}

	if !userAccount(user.AccountIDs, req.AccountID) {
		return BalanceResponse{}, errors.New("unauthorized account id")
	}

	account, err := s.accountRepo.Read(
		accountrepo.ReadRequest{
			ID: req.AccountID,
		},
	)

	if err != nil {
		return BalanceResponse{}, err
	}

	return BalanceResponse{
		Balance: account.Balance,
	}, nil
}

func (s svc) History(req HistoryRequest) (HistoryResponse, error) {
	if !validID(req.UserID) {
		return HistoryResponse{}, errors.New("invalid user id")
	}

	if !validID(req.AccountID) {
		return HistoryResponse{}, errors.New("invalid account id")
	}

	user, err := s.userRepo.Read(
		userrepo.ReadRequest{
			ID: req.UserID,
		},
	)

	if err != nil {
		return HistoryResponse{}, err
	}

	if !userAccount(user.AccountIDs, req.AccountID) {
		return HistoryResponse{}, errors.New("unauthorized account id")
	}

	account, err := s.accountRepo.Read(
		accountrepo.ReadRequest{
			ID: req.AccountID,
		},
	)

	if err != nil {
		return HistoryResponse{}, err
	}

	history, err := s.historyRepo.Read(
		historyrepo.ReadRequest{
			ID: account.HistoryID,
		},
	)

	if err != nil {
		return HistoryResponse{}, err
	}

	return HistoryResponse{
		Events: history.Events,
	}, nil
}

func validID(id string) bool {
	if id == "" {
		return false
	}

	_, err := guuid.Parse(id)

	return err == nil
}

func userAccount(userAccountIDs []string, accountID string) bool {
	for _, userAccountID := range userAccountIDs {
		if userAccountID == accountID {
			return true
		}
	}

	return false
}
