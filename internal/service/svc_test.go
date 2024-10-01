package service

import (
	"errors"
	"testing"
	"time"

	"github.com/hetfdex/tiny-bank/internal/domain"
	"github.com/hetfdex/tiny-bank/internal/repository/accountrepo"
	"github.com/hetfdex/tiny-bank/internal/repository/userrepo"
	"github.com/hetfdex/tiny-bank/test/mock/repository/accountrepomock"
	"github.com/hetfdex/tiny-bank/test/mock/repository/historyrepomock"
	"github.com/hetfdex/tiny-bank/test/mock/repository/userrepomock"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransfer_ErrInvalidSenderUserID(t *testing.T) {
	svc := New(nil, nil, nil)

	res, err := svc.Transfer(TransferRequest{})

	assert.Equal(t, TransferResponse{}, res)
	assert.Equal(t, errors.New("invalid sender user id"), err)
}

func TestTransfer_ErrInvalidReceiverUserID(t *testing.T) {
	svc := New(nil, nil, nil)

	res, err := svc.Transfer(
		TransferRequest{
			SenderUserID:   uuid.New(),
			ReceiverUserID: "1234",
		},
	)

	assert.Equal(t, TransferResponse{}, res)
	assert.Equal(t, errors.New("invalid receiver user id"), err)
}

func TestTransfer_ErrInvalidSenderAccountID(t *testing.T) {
	svc := New(nil, nil, nil)

	res, err := svc.Transfer(
		TransferRequest{
			SenderUserID:   uuid.New(),
			ReceiverUserID: uuid.New(),
		},
	)

	assert.Equal(t, TransferResponse{}, res)
	assert.Equal(t, errors.New("invalid sender account id"), err)
}

func TestTransfer_ErrInvalidReceiverAccountID(t *testing.T) {
	svc := New(nil, nil, nil)

	res, err := svc.Transfer(
		TransferRequest{
			SenderUserID:      uuid.New(),
			ReceiverUserID:    uuid.New(),
			SenderAccountID:   uuid.New(),
			ReceiverAccountID: "1234",
		},
	)

	assert.Equal(t, TransferResponse{}, res)
	assert.Equal(t, errors.New("invalid receiver account id"), err)
}

func TestTransfer_ErrInvalidAmount(t *testing.T) {
	svc := New(nil, nil, nil)

	userID := uuid.New()
	accountID := uuid.New()

	res, err := svc.Transfer(
		TransferRequest{
			SenderUserID:      userID,
			ReceiverUserID:    userID,
			SenderAccountID:   accountID,
			ReceiverAccountID: accountID,
			Amount:            -1,
		},
	)

	assert.Equal(t, TransferResponse{}, res)
	assert.Equal(t, errors.New("invalid amount"), err)
}

func TestTransfer_ErrSameAccount(t *testing.T) {
	svc := New(nil, nil, nil)

	userID := uuid.New()
	accountID := uuid.New()

	res, err := svc.Transfer(
		TransferRequest{
			SenderUserID:      userID,
			ReceiverUserID:    userID,
			SenderAccountID:   accountID,
			ReceiverAccountID: accountID,
			Amount:            1,
		},
	)

	assert.Equal(t, TransferResponse{}, res)
	assert.Equal(t, errors.New("same account"), err)
}

func TestTransfer_ErrReadSender(t *testing.T) {
	errMock := errors.New("error user")

	senderUserID := uuid.New()

	userRepo := &userrepomock.Mock{}

	userRepo.On(
		"Read",
		userrepo.ReadRequest{
			ID: senderUserID,
		},
	).Return(
		domain.User{},
		errMock,
	)

	svc := New(userRepo, nil, nil)

	res, err := svc.Transfer(
		TransferRequest{
			SenderUserID:      senderUserID,
			ReceiverUserID:    uuid.New(),
			SenderAccountID:   uuid.New(),
			ReceiverAccountID: uuid.New(),
			Amount:            10,
		},
	)

	assert.Equal(t, TransferResponse{}, res)
	assert.Equal(t, errMock, err)
}

func TestTransfer_ErrUnauthorizedAccountIDSender(t *testing.T) {
	senderUserID := uuid.New()

	userRepo := &userrepomock.Mock{}

	userRepo.On(
		"Read",
		userrepo.ReadRequest{
			ID: senderUserID,
		},
	).Return(
		domain.User{
			ID:         senderUserID,
			CreatedAt:  time.Now().UTC(),
			Active:     true,
			Name:       "joe",
			AccountIDs: []string{uuid.New()},
		},
		nil,
	)

	svc := New(userRepo, nil, nil)

	res, err := svc.Transfer(
		TransferRequest{
			SenderUserID:      senderUserID,
			ReceiverUserID:    uuid.New(),
			SenderAccountID:   uuid.New(),
			ReceiverAccountID: uuid.New(),
			Amount:            10,
		},
	)

	assert.Equal(t, TransferResponse{}, res)
	assert.Equal(t, errors.New("unauthorized account id"), err)
}

func TestTransfer_ErrReadSenderAccount(t *testing.T) {
	errMock := errors.New("error account")

	senderUserID := uuid.New()
	senderAccountID := uuid.New()

	userRepo := &userrepomock.Mock{}

	userRepo.On(
		"Read",
		userrepo.ReadRequest{
			ID: senderUserID,
		},
	).Return(
		domain.User{
			ID:         senderUserID,
			CreatedAt:  time.Now().UTC(),
			Active:     true,
			Name:       "joe",
			AccountIDs: []string{senderAccountID},
		},
		nil,
	)

	accountRepo := &accountrepomock.Mock{}

	accountRepo.On(
		"Read",
		accountrepo.ReadRequest{
			ID: senderAccountID,
		},
	).Return(
		domain.Account{},
		errMock,
	)

	svc := New(userRepo, accountRepo, nil)

	res, err := svc.Transfer(
		TransferRequest{
			SenderUserID:      senderUserID,
			ReceiverUserID:    uuid.New(),
			SenderAccountID:   senderAccountID,
			ReceiverAccountID: uuid.New(),
			Amount:            10,
		},
	)

	assert.Equal(t, TransferResponse{}, res)
	assert.Equal(t, errMock, err)
}

func TestTransfer_ErrInsuficientFunds(t *testing.T) {
	senderUserID := uuid.New()
	senderAccountID := uuid.New()

	userRepo := &userrepomock.Mock{}

	userRepo.On(
		"Read",
		userrepo.ReadRequest{
			ID: senderUserID,
		},
	).Return(
		domain.User{
			ID:         senderUserID,
			CreatedAt:  time.Now().UTC(),
			Active:     true,
			Name:       "joe",
			AccountIDs: []string{senderAccountID},
		},
		nil,
	)

	accountRepo := &accountrepomock.Mock{}

	accountRepo.On(
		"Read",
		accountrepo.ReadRequest{
			ID: senderAccountID,
		},
	).Return(
		domain.Account{
			ID:        senderAccountID,
			CreatedAt: time.Now().UTC(),
			Balance:   0,
			HistoryID: uuid.New(),
		},
		nil,
	)

	svc := New(userRepo, accountRepo, nil)

	res, err := svc.Transfer(
		TransferRequest{
			SenderUserID:      senderUserID,
			ReceiverUserID:    uuid.New(),
			SenderAccountID:   senderAccountID,
			ReceiverAccountID: uuid.New(),
			Amount:            10,
		},
	)

	assert.Equal(t, TransferResponse{}, res)
	assert.Equal(t, errors.New("insuficient funds"), err)
}

func TestTransfer_ErrReadReceiver(t *testing.T) {
	errMock := errors.New("error user")

	senderUserID := uuid.New()
	receiverUserID := uuid.New()
	senderAccountID := uuid.New()

	userRepo := &userrepomock.Mock{}

	userRepo.On(
		"Read",
		userrepo.ReadRequest{
			ID: senderUserID,
		},
	).Return(
		domain.User{
			ID:         senderUserID,
			CreatedAt:  time.Now().UTC(),
			Active:     true,
			Name:       "joe",
			AccountIDs: []string{senderAccountID},
		},
		nil,
	)

	userRepo.On(
		"Read",
		userrepo.ReadRequest{
			ID: receiverUserID,
		},
	).Return(
		domain.User{},
		errMock,
	)

	accountRepo := &accountrepomock.Mock{}

	accountRepo.On(
		"Read",
		accountrepo.ReadRequest{
			ID: senderAccountID,
		},
	).Return(
		domain.Account{
			ID:        senderAccountID,
			CreatedAt: time.Now().UTC(),
			Balance:   20,
			HistoryID: uuid.New(),
		},
		nil,
	)

	svc := New(userRepo, accountRepo, nil)

	res, err := svc.Transfer(
		TransferRequest{
			SenderUserID:      senderUserID,
			ReceiverUserID:    receiverUserID,
			SenderAccountID:   senderAccountID,
			ReceiverAccountID: uuid.New(),
			Amount:            10,
		},
	)

	assert.Equal(t, TransferResponse{}, res)
	assert.Equal(t, errMock, err)
}

func TestTransfer_ErrUnauthorizedAccountIDReceiver(t *testing.T) {
	senderUserID := uuid.New()
	receiverUserID := uuid.New()
	senderAccountID := uuid.New()
	receiverAccountID := uuid.New()

	userRepo := &userrepomock.Mock{}

	userRepo.On(
		"Read",
		userrepo.ReadRequest{
			ID: senderUserID,
		},
	).Return(
		domain.User{
			ID:         senderUserID,
			CreatedAt:  time.Now().UTC(),
			Active:     true,
			Name:       "joe",
			AccountIDs: []string{senderAccountID},
		},
		nil,
	)

	userRepo.On(
		"Read",
		userrepo.ReadRequest{
			ID: receiverUserID,
		},
	).Return(
		domain.User{
			ID:         receiverUserID,
			CreatedAt:  time.Now().UTC(),
			Active:     true,
			Name:       "mary",
			AccountIDs: []string{uuid.New()},
		},
		nil,
	)

	accountRepo := &accountrepomock.Mock{}

	accountRepo.On(
		"Read",
		accountrepo.ReadRequest{
			ID: senderAccountID,
		},
	).Return(
		domain.Account{
			ID:        senderAccountID,
			CreatedAt: time.Now().UTC(),
			Balance:   20,
			HistoryID: uuid.New(),
		},
		nil,
	)

	svc := New(userRepo, accountRepo, nil)

	res, err := svc.Transfer(
		TransferRequest{
			SenderUserID:      senderUserID,
			ReceiverUserID:    receiverUserID,
			SenderAccountID:   senderAccountID,
			ReceiverAccountID: receiverAccountID,
			Amount:            10,
		},
	)

	assert.Equal(t, TransferResponse{}, res)
	assert.Equal(t, errors.New("unauthorized account id"), err)
}

func TestTransfer_ErrReadReceiverAccount(t *testing.T) {
	errMock := errors.New("error account")

	senderUserID := uuid.New()
	receiverUserID := uuid.New()
	senderAccountID := uuid.New()
	receiverAccountID := uuid.New()

	userRepo := &userrepomock.Mock{}

	userRepo.On(
		"Read",
		userrepo.ReadRequest{
			ID: senderUserID,
		},
	).Return(
		domain.User{
			ID:         senderUserID,
			CreatedAt:  time.Now().UTC(),
			Active:     true,
			Name:       "joe",
			AccountIDs: []string{senderAccountID},
		},
		nil,
	)

	userRepo.On(
		"Read",
		userrepo.ReadRequest{
			ID: receiverUserID,
		},
	).Return(
		domain.User{
			ID:         receiverUserID,
			CreatedAt:  time.Now().UTC(),
			Active:     true,
			Name:       "mary",
			AccountIDs: []string{receiverAccountID},
		},
		nil,
	)

	accountRepo := &accountrepomock.Mock{}

	accountRepo.On(
		"Read",
		accountrepo.ReadRequest{
			ID: senderAccountID,
		},
	).Return(
		domain.Account{
			ID:        senderAccountID,
			CreatedAt: time.Now().UTC(),
			Balance:   20,
			HistoryID: uuid.New(),
		},
		nil,
	)

	accountRepo.On(
		"Read",
		accountrepo.ReadRequest{
			ID: receiverAccountID,
		},
	).Return(
		domain.Account{},
		errMock,
	)

	svc := New(userRepo, accountRepo, nil)

	res, err := svc.Transfer(
		TransferRequest{
			SenderUserID:      senderUserID,
			ReceiverUserID:    receiverUserID,
			SenderAccountID:   senderAccountID,
			ReceiverAccountID: receiverAccountID,
			Amount:            10,
		},
	)

	assert.Equal(t, TransferResponse{}, res)
	assert.Equal(t, errMock, err)
}

func TestTransfer_ErrUpdateSenderAccount(t *testing.T) {
	errMock := errors.New("error account")

	senderUserID := uuid.New()
	receiverUserID := uuid.New()
	senderAccountID := uuid.New()
	receiverAccountID := uuid.New()

	userRepo := &userrepomock.Mock{}

	userRepo.On(
		"Read",
		userrepo.ReadRequest{
			ID: senderUserID,
		},
	).Return(
		domain.User{
			ID:         senderUserID,
			CreatedAt:  time.Now().UTC(),
			Active:     true,
			Name:       "joe",
			AccountIDs: []string{senderAccountID},
		},
		nil,
	)

	userRepo.On(
		"Read",
		userrepo.ReadRequest{
			ID: receiverUserID,
		},
	).Return(
		domain.User{
			ID:         receiverUserID,
			CreatedAt:  time.Now().UTC(),
			Active:     true,
			Name:       "mary",
			AccountIDs: []string{receiverAccountID},
		},
		nil,
	)

	accountRepo := &accountrepomock.Mock{}

	accountRepo.On(
		"Read",
		accountrepo.ReadRequest{
			ID: senderAccountID,
		},
	).Return(
		domain.Account{
			ID:        senderAccountID,
			CreatedAt: time.Now().UTC(),
			Balance:   20,
			HistoryID: uuid.New(),
		},
		nil,
	)

	accountRepo.On(
		"Read",
		accountrepo.ReadRequest{
			ID: receiverAccountID,
		},
	).Return(
		domain.Account{
			ID:        receiverAccountID,
			CreatedAt: time.Now().UTC(),
			Balance:   10,
			HistoryID: uuid.New(),
		},
		nil,
	)

	accountRepo.On(
		"Update",
		accountrepo.UpdateRequest{
			ID:      senderAccountID,
			Balance: 10,
		},
	).Return(
		errMock,
	)

	svc := New(userRepo, accountRepo, nil)

	res, err := svc.Transfer(
		TransferRequest{
			SenderUserID:      senderUserID,
			ReceiverUserID:    receiverUserID,
			SenderAccountID:   senderAccountID,
			ReceiverAccountID: receiverAccountID,
			Amount:            10,
		},
	)

	assert.Equal(t, TransferResponse{}, res)
	assert.Equal(t, errMock, err)
}

func TestTransfer_ErrUpdateReceiverAccount(t *testing.T) {
	errMock := errors.New("error account")

	senderUserID := uuid.New()
	receiverUserID := uuid.New()
	senderAccountID := uuid.New()
	receiverAccountID := uuid.New()

	userRepo := &userrepomock.Mock{}

	userRepo.On(
		"Read",
		userrepo.ReadRequest{
			ID: senderUserID,
		},
	).Return(
		domain.User{
			ID:         senderUserID,
			CreatedAt:  time.Now().UTC(),
			Active:     true,
			Name:       "joe",
			AccountIDs: []string{senderAccountID},
		},
		nil,
	)

	userRepo.On(
		"Read",
		userrepo.ReadRequest{
			ID: receiverUserID,
		},
	).Return(
		domain.User{
			ID:         receiverUserID,
			CreatedAt:  time.Now().UTC(),
			Active:     true,
			Name:       "mary",
			AccountIDs: []string{receiverAccountID},
		},
		nil,
	)

	accountRepo := &accountrepomock.Mock{}

	accountRepo.On(
		"Read",
		accountrepo.ReadRequest{
			ID: senderAccountID,
		},
	).Return(
		domain.Account{
			ID:        senderAccountID,
			CreatedAt: time.Now().UTC(),
			Balance:   20,
			HistoryID: uuid.New(),
		},
		nil,
	)

	accountRepo.On(
		"Read",
		accountrepo.ReadRequest{
			ID: receiverAccountID,
		},
	).Return(
		domain.Account{
			ID:        receiverAccountID,
			CreatedAt: time.Now().UTC(),
			Balance:   10,
			HistoryID: uuid.New(),
		},
		nil,
	)

	accountRepo.On(
		"Update",
		accountrepo.UpdateRequest{
			ID:      senderAccountID,
			Balance: 10,
		},
	).Return(
		nil,
	)

	accountRepo.On(
		"Update",
		accountrepo.UpdateRequest{
			ID:      receiverAccountID,
			Balance: 20,
		},
	).Return(
		errMock,
	)

	svc := New(userRepo, accountRepo, nil)

	res, err := svc.Transfer(
		TransferRequest{
			SenderUserID:      senderUserID,
			ReceiverUserID:    receiverUserID,
			SenderAccountID:   senderAccountID,
			ReceiverAccountID: receiverAccountID,
			Amount:            10,
		},
	)

	assert.Equal(t, TransferResponse{}, res)
	assert.Equal(t, errMock, err)
}

func TestTransfer_ErrUpdateSenderHistory(t *testing.T) {
	errMock := errors.New("error history")

	senderUserID := uuid.New()
	receiverUserID := uuid.New()
	senderAccountID := uuid.New()
	receiverAccountID := uuid.New()

	userRepo := &userrepomock.Mock{}

	userRepo.On(
		"Read",
		userrepo.ReadRequest{
			ID: senderUserID,
		},
	).Return(
		domain.User{
			ID:         senderUserID,
			CreatedAt:  time.Now().UTC(),
			Active:     true,
			Name:       "joe",
			AccountIDs: []string{senderAccountID},
		},
		nil,
	)

	userRepo.On(
		"Read",
		userrepo.ReadRequest{
			ID: receiverUserID,
		},
	).Return(
		domain.User{
			ID:         receiverUserID,
			CreatedAt:  time.Now().UTC(),
			Active:     true,
			Name:       "mary",
			AccountIDs: []string{receiverAccountID},
		},
		nil,
	)

	accountRepo := &accountrepomock.Mock{}

	accountRepo.On(
		"Read",
		accountrepo.ReadRequest{
			ID: senderAccountID,
		},
	).Return(
		domain.Account{
			ID:        senderAccountID,
			CreatedAt: time.Now().UTC(),
			Balance:   20,
			HistoryID: uuid.New(),
		},
		nil,
	)

	accountRepo.On(
		"Read",
		accountrepo.ReadRequest{
			ID: receiverAccountID,
		},
	).Return(
		domain.Account{
			ID:        receiverAccountID,
			CreatedAt: time.Now().UTC(),
			Balance:   10,
			HistoryID: uuid.New(),
		},
		nil,
	)

	accountRepo.On(
		"Update",
		accountrepo.UpdateRequest{
			ID:      senderAccountID,
			Balance: 10,
		},
	).Return(
		nil,
	)

	accountRepo.On(
		"Update",
		accountrepo.UpdateRequest{
			ID:      receiverAccountID,
			Balance: 20,
		},
	).Return(
		nil,
	)

	historyRepo := &historyrepomock.Mock{}

	historyRepo.On(
		"Update",
		mock.AnythingOfType("historyrepo.UpdateRequest"),
	).Return(
		errMock,
	)

	svc := New(userRepo, accountRepo, historyRepo)

	res, err := svc.Transfer(
		TransferRequest{
			SenderUserID:      senderUserID,
			ReceiverUserID:    receiverUserID,
			SenderAccountID:   senderAccountID,
			ReceiverAccountID: receiverAccountID,
			Amount:            10,
		},
	)

	assert.Equal(t, TransferResponse{}, res)
	assert.Equal(t, errMock, err)
}

func TestTransfer_ErrUpdateReceiverHistory(t *testing.T) {
	errMock := errors.New("error history")

	senderUserID := uuid.New()
	receiverUserID := uuid.New()
	senderAccountID := uuid.New()
	receiverAccountID := uuid.New()

	userRepo := &userrepomock.Mock{}

	userRepo.On(
		"Read",
		userrepo.ReadRequest{
			ID: senderUserID,
		},
	).Return(
		domain.User{
			ID:         senderUserID,
			CreatedAt:  time.Now().UTC(),
			Active:     true,
			Name:       "joe",
			AccountIDs: []string{senderAccountID},
		},
		nil,
	)

	userRepo.On(
		"Read",
		userrepo.ReadRequest{
			ID: receiverUserID,
		},
	).Return(
		domain.User{
			ID:         receiverUserID,
			CreatedAt:  time.Now().UTC(),
			Active:     true,
			Name:       "mary",
			AccountIDs: []string{receiverAccountID},
		},
		nil,
	)

	accountRepo := &accountrepomock.Mock{}

	accountRepo.On(
		"Read",
		accountrepo.ReadRequest{
			ID: senderAccountID,
		},
	).Return(
		domain.Account{
			ID:        senderAccountID,
			CreatedAt: time.Now().UTC(),
			Balance:   20,
			HistoryID: uuid.New(),
		},
		nil,
	)

	accountRepo.On(
		"Read",
		accountrepo.ReadRequest{
			ID: receiverAccountID,
		},
	).Return(
		domain.Account{
			ID:        receiverAccountID,
			CreatedAt: time.Now().UTC(),
			Balance:   10,
			HistoryID: uuid.New(),
		},
		nil,
	)

	accountRepo.On(
		"Update",
		accountrepo.UpdateRequest{
			ID:      senderAccountID,
			Balance: 10,
		},
	).Return(
		nil,
	)

	accountRepo.On(
		"Update",
		accountrepo.UpdateRequest{
			ID:      receiverAccountID,
			Balance: 20,
		},
	).Return(
		nil,
	)

	historyRepo := &historyrepomock.Mock{}

	historyRepo.On(
		"Update",
		mock.AnythingOfType("historyrepo.UpdateRequest"),
	).Return(
		nil,
	).Once()

	historyRepo.On(
		"Update",
		mock.AnythingOfType("historyrepo.UpdateRequest"),
	).Return(
		errMock,
	).Once()

	svc := New(userRepo, accountRepo, historyRepo)

	res, err := svc.Transfer(
		TransferRequest{
			SenderUserID:      senderUserID,
			ReceiverUserID:    receiverUserID,
			SenderAccountID:   senderAccountID,
			ReceiverAccountID: receiverAccountID,
			Amount:            10,
		},
	)

	assert.Equal(t, TransferResponse{}, res)
	assert.Equal(t, errMock, err)
}

func TestTransfer_Ok(t *testing.T) {
	senderUserID := uuid.New()
	receiverUserID := uuid.New()
	senderAccountID := uuid.New()
	receiverAccountID := uuid.New()

	userRepo := &userrepomock.Mock{}

	userRepo.On(
		"Read",
		userrepo.ReadRequest{
			ID: senderUserID,
		},
	).Return(
		domain.User{
			ID:         senderUserID,
			CreatedAt:  time.Now().UTC(),
			Active:     true,
			Name:       "joe",
			AccountIDs: []string{senderAccountID},
		},
		nil,
	)

	userRepo.On(
		"Read",
		userrepo.ReadRequest{
			ID: receiverUserID,
		},
	).Return(
		domain.User{
			ID:         receiverUserID,
			CreatedAt:  time.Now().UTC(),
			Active:     true,
			Name:       "mary",
			AccountIDs: []string{receiverAccountID},
		},
		nil,
	)

	accountRepo := &accountrepomock.Mock{}

	accountRepo.On(
		"Read",
		accountrepo.ReadRequest{
			ID: senderAccountID,
		},
	).Return(
		domain.Account{
			ID:        senderAccountID,
			CreatedAt: time.Now().UTC(),
			Balance:   20,
			HistoryID: uuid.New(),
		},
		nil,
	)

	accountRepo.On(
		"Read",
		accountrepo.ReadRequest{
			ID: receiverAccountID,
		},
	).Return(
		domain.Account{
			ID:        receiverAccountID,
			CreatedAt: time.Now().UTC(),
			Balance:   10,
			HistoryID: uuid.New(),
		},
		nil,
	)

	accountRepo.On(
		"Update",
		accountrepo.UpdateRequest{
			ID:      senderAccountID,
			Balance: 10,
		},
	).Return(
		nil,
	)

	accountRepo.On(
		"Update",
		accountrepo.UpdateRequest{
			ID:      receiverAccountID,
			Balance: 20,
		},
	).Return(
		nil,
	)

	historyRepo := &historyrepomock.Mock{}

	historyRepo.On(
		"Update",
		mock.AnythingOfType("historyrepo.UpdateRequest"),
	).Return(
		nil,
	)

	svc := New(userRepo, accountRepo, historyRepo)

	res, err := svc.Transfer(
		TransferRequest{
			SenderUserID:      senderUserID,
			ReceiverUserID:    receiverUserID,
			SenderAccountID:   senderAccountID,
			ReceiverAccountID: receiverAccountID,
			Amount:            10,
		},
	)

	assert.Equal(t, TransferResponse{Balance: 10}, res)
	assert.Nil(t, err)
}
