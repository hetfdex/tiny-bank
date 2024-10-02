package userrepo

import (
	"errors"
	"sync"
	"time"

	"github.com/hetfdex/tiny-bank/internal/domain"
	"github.com/pborman/uuid"
)

var (
	usersMux sync.Mutex
)

type Repo interface {
	Create(CreateRequest) (domain.User, error)
	Read(ReadRequest) (domain.User, error)
	UpdateStatus(UpdateStatusRequest) error
	UpdateAccountIDs(UpdateAccountIDsRequest) error
}

type repo struct {
	users map[string]domain.User
}

func New(
	users map[string]domain.User,
) Repo {

	return &repo{
		users: users,
	}
}

func (r repo) Create(req CreateRequest) (domain.User, error) {
	usersMux.Lock()

	defer usersMux.Unlock()

	id := uuid.New()

	if _, exists := r.users[id]; exists {
		return domain.User{}, errors.New("duplicate user id")
	}

	user := domain.User{
		ID:         id,
		CreatedAt:  time.Now().UTC(),
		Active:     true,
		Name:       req.Name,
		AccountIDs: map[string]struct{}{},
	}

	r.users[id] = user

	return user, nil
}

func (r repo) Read(req ReadRequest) (domain.User, error) {
	usersMux.Lock()

	defer usersMux.Unlock()

	return r.getActiveUser(req.ID)
}

func (r repo) UpdateStatus(req UpdateStatusRequest) error {
	usersMux.Lock()

	defer usersMux.Unlock()

	user, exists := r.users[req.ID]

	if !exists {
		return errors.New("user not found")
	}

	user.Active = req.Active

	r.users[req.ID] = user

	return nil
}

func (r repo) UpdateAccountIDs(req UpdateAccountIDsRequest) error {
	usersMux.Lock()

	defer usersMux.Unlock()

	user, err := r.getActiveUser(req.ID)

	if err != nil {
		return err
	}

	if _, exists := user.AccountIDs[req.AccountID]; exists {
		return errors.New("duplicate account id")
	}

	user.AccountIDs[req.AccountID] = struct{}{}

	r.users[req.ID] = user

	return nil
}

func (r repo) getActiveUser(id string) (domain.User, error) {
	user, exists := r.users[id]

	if !exists {
		return domain.User{}, errors.New("user not found")
	}

	if !user.Active {
		return domain.User{}, errors.New("user not active")
	}

	return user, nil
}
