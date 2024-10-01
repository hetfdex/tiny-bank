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
	Update(UpdateRequest) error
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
		Balance:   0,
		HistoryID: req.HistoryID,
	}

	r.accounts[id] = account

	return account, nil
}

func (r repo) Read(req ReadRequest) (domain.Account, error) {
	accountsMux.Lock()

	defer accountsMux.Unlock()

	account, exists := r.accounts[req.ID]

	if !exists {
		return domain.Account{}, errors.New("account not found")
	}

	return account, nil
}

func (r repo) Update(req UpdateRequest) error {
	accountsMux.Lock()

	defer accountsMux.Unlock()

	account, exists := r.accounts[req.ID]

	if !exists {
		return errors.New("account not found")
	}

	if account.Balance == req.Balance {
		return nil
	}

	account.Balance = int(req.Balance)

	r.accounts[req.ID] = account

	return nil
}
