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
	Update(UpdateRequest) error
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
		return domain.User{}, errors.New("id in use")
	}

	user := domain.User{
		ID:         id,
		CreatedAt:  time.Now().UTC(),
		Active:     true,
		Name:       req.Name,
		AccountIDs: []string{req.AccountID},
	}

	r.users[id] = user

	return user, nil
}

func (r repo) Read(req ReadRequest) (domain.User, error) {
	usersMux.Lock()

	defer usersMux.Unlock()

	user, exists := r.users[req.ID]

	if !exists {
		return domain.User{}, errors.New("user not found")
	}

	if !user.Active {
		return domain.User{}, errors.New("user not active")
	}

	return user, nil
}

func (r repo) Update(req UpdateRequest) error {
	usersMux.Lock()

	defer usersMux.Unlock()

	user, exists := r.users[req.ID]

	if !exists {
		return errors.New("user not found")
	}

	if user.Active == req.Active {
		return nil
	}

	user.Active = req.Active

	r.users[req.ID] = user

	return nil
}
