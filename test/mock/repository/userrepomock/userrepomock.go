package userrepomock

import (
	"github.com/hetfdex/tiny-bank/internal/domain"
	"github.com/hetfdex/tiny-bank/internal/repository/userrepo"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) Create(req userrepo.CreateRequest) (domain.User, error) {
	args := m.Called(req)

	return args.Get(0).(domain.User), args.Error(1)
}

func (m *Mock) Read(req userrepo.ReadRequest) (domain.User, error) {
	args := m.Called(req)

	return args.Get(0).(domain.User), args.Error(1)
}

func (m *Mock) Update(req userrepo.UpdateRequest) error {
	args := m.Called(req)

	return args.Error(0)
}
