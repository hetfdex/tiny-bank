package accountrepomock

import (
	"github.com/hetfdex/tiny-bank/internal/domain"
	"github.com/hetfdex/tiny-bank/internal/repository/accountrepo"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) Create(req accountrepo.CreateRequest) (domain.Account, error) {
	args := m.Called(req)

	return args.Get(0).(domain.Account), args.Error(1)
}

func (m *Mock) Read(req accountrepo.ReadRequest) (domain.Account, error) {
	args := m.Called(req)

	return args.Get(0).(domain.Account), args.Error(1)
}

func (m *Mock) Update(req accountrepo.UpdateRequest) error {
	args := m.Called(req)

	return args.Error(0)
}
