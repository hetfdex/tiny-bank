package historyrepo

import (
	"errors"
	"sync"
	"time"

	"github.com/hetfdex/tiny-bank/internal/domain"
	"github.com/pborman/uuid"
)

var (
	historiesMux sync.Mutex
)

type Repo interface {
	Create(CreateRequest) (domain.History, error)
	Read(ReadRequest) (domain.History, error)
	Update(UpdateRequest) error
}

type repo struct {
	histories map[string]domain.History
}

func New(
	histories map[string]domain.History,
) Repo {

	return &repo{
		histories: histories,
	}
}
func (r repo) Create(req CreateRequest) (domain.History, error) {
	historiesMux.Lock()

	defer historiesMux.Unlock()

	id := uuid.New()

	if _, exists := r.histories[id]; exists {
		return domain.History{}, errors.New("id in use")
	}

	history := domain.History{
		ID:        id,
		CreatedAt: time.Now().UTC(),
		Events:    []domain.Event{},
	}

	r.histories[id] = history

	return history, nil
}

func (r repo) Read(req ReadRequest) (domain.History, error) {
	historiesMux.Lock()

	defer historiesMux.Unlock()

	history, exists := r.histories[req.ID]

	if !exists {
		return domain.History{}, errors.New("history not found")
	}

	return history, nil
}

func (r repo) Update(req UpdateRequest) error {
	historiesMux.Lock()

	defer historiesMux.Unlock()

	history, exists := r.histories[req.ID]

	if !exists {
		return errors.New("history not found")
	}

	if req.Event == (domain.Event{}) {
		return nil
	}

	history.Events = append(history.Events, req.Event)

	r.histories[req.ID] = history

	return nil
}
