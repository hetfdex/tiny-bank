package historyrepomock

import (
	"github.com/hetfdex/tiny-bank/internal/domain"
	"github.com/hetfdex/tiny-bank/internal/repository/historyrepo"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) Create(req historyrepo.CreateRequest) (domain.History, error) {
	args := m.Called(req)

	return args.Get(0).(domain.History), args.Error(1)
}

func (m *Mock) Read(req historyrepo.ReadRequest) (domain.History, error) {
	args := m.Called(req)

	return args.Get(0).(domain.History), args.Error(1)
}

func (m *Mock) Update(req historyrepo.UpdateRequest) error {
	args := m.Called(req)

	return args.Error(0)
}
