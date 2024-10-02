package accountrepo

import (
	"errors"
	"sync"
	"time"

	"github.com/hetfdex/tiny-bank/internal/domain"
	"github.com/pborman/uuid"
)

var (
	accountsMux sync.Mutex
)

type Repo interface {
	Create(CreateRequest) (domain.Account, error)
	Read(ReadRequest) (domain.Account, error)
	UpdateBalance(UpdateBalanceRequest) error
	UpdateTransactions(UpdateTransactionsRequest) error
}

type repo struct {
	accounts map[string]domain.Account
}

func New(
	accounts map[string]domain.Account,
) Repo {

	return &repo{
		accounts: accounts,
	}
}

func (r repo) Create(req CreateRequest) (domain.Account, error) {
	accountsMux.Lock()

	defer accountsMux.Unlock()

	id := uuid.New()

	if _, exists := r.accounts[id]; exists {
		return domain.Account{}, errors.New("id in use")
	}

	account := domain.Account{
		ID:        id,
		CreatedAt: time.Now().UTC(),
	}

	r.accounts[id] = account

	return account, nil
}

func (r repo) Read(req ReadRequest) (domain.Account, error) {
	accountsMux.Lock()

	defer accountsMux.Unlock()

	return r.getAccount(req.ID)
}

func (r repo) UpdateBalance(req UpdateBalanceRequest) error {
	accountsMux.Lock()

	defer accountsMux.Unlock()

	account, err := r.getAccount(req.ID)

	if err != nil {
		return err
	}

	account.Balance = int(req.Balance)

	r.accounts[req.ID] = account

	return nil
}

func (r repo) UpdateTransactions(req UpdateTransactionsRequest) error {
	accountsMux.Lock()

	defer accountsMux.Unlock()

	account, err := r.getAccount(req.ID)

	if err != nil {
		return err
	}

	account.Transactions = append(account.Transactions, req.Transaction)

	r.accounts[req.ID] = account

	return nil
}

func (r repo) getAccount(id string) (domain.Account, error) {
	account, exists := r.accounts[id]

	if !exists {
		return domain.Account{}, errors.New("account not found")
	}

	return account, nil
}
