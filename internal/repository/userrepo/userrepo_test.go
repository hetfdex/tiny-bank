package userrepo

import (
	"errors"
	"testing"
	"time"

	"github.com/hetfdex/tiny-bank/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreate_Ok(t *testing.T) {
	repo := New(make(map[string]domain.User))

	res, err := repo.Create(
		CreateRequest{
			Name:      "joe",
			AccountID: "1234",
		},
	)

	assert.NotEmpty(t, res.ID)
	assert.NotEmpty(t, res.CreatedAt)
	assert.True(t, res.Active)
	assert.Equal(t, "joe", res.Name)
	assert.Equal(t, 1, len(res.AccountIDs))
	assert.Equal(t, "1234", res.AccountIDs[0])

	assert.Nil(t, err)
}

func TestRead_ErrUserNotFound(t *testing.T) {
	repo := New(make(map[string]domain.User))

	res, err := repo.Read(
		ReadRequest{
			ID: "1234",
		},
	)

	assert.Equal(t, domain.User{}, res)
	assert.Equal(t, errors.New("user not found"), err)
}

func TestRead_ErrUserNotActive(t *testing.T) {
	user := domain.User{
		ID:         "1234",
		CreatedAt:  time.Now().UTC(),
		Active:     false,
		Name:       "joe",
		AccountIDs: []string{"5678"},
	}

	users := make(map[string]domain.User)

	users["1234"] = user

	repo := New(users)

	res, err := repo.Read(
		ReadRequest{
			ID: "1234",
		},
	)

	assert.Equal(t, domain.User{}, res)
	assert.Equal(t, errors.New("user not active"), err)
}

func TestRead_Ok(t *testing.T) {
	user := domain.User{
		ID:         "1234",
		CreatedAt:  time.Now().UTC(),
		Active:     true,
		Name:       "joe",
		AccountIDs: []string{"5678"},
	}

	users := make(map[string]domain.User)

	users["1234"] = user

	repo := New(users)

	res, err := repo.Read(
		ReadRequest{
			ID: "1234",
		},
	)

	assert.Equal(t, user, res)
	assert.Nil(t, err)
}

func TestUpdate_ErrUserNotFound(t *testing.T) {
	repo := New(make(map[string]domain.User))

	err := repo.Update(
		UpdateRequest{
			ID: "1234",
		},
	)

	assert.Equal(t, errors.New("user not found"), err)
}

func TestUpdate_OkNoChange(t *testing.T) {
	users := make(map[string]domain.User)

	users["1234"] = domain.User{
		ID:         "1234",
		CreatedAt:  time.Now().UTC(),
		Active:     false,
		Name:       "joe",
		AccountIDs: []string{"5678"},
	}

	repo := New(users)

	err := repo.Update(
		UpdateRequest{
			ID:     "1234",
			Active: false,
		},
	)

	assert.Nil(t, err)
}

func TestUpdate_Ok(t *testing.T) {
	users := make(map[string]domain.User)

	users["1234"] = domain.User{
		ID:         "1234",
		CreatedAt:  time.Now().UTC(),
		Active:     true,
		Name:       "joe",
		AccountIDs: []string{"5678"},
	}

	repo := New(users)

	err := repo.Update(
		UpdateRequest{
			ID:     "1234",
			Active: false,
		},
	)

	assert.Nil(t, err)
}
