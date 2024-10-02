package accountrepo

import "github.com/hetfdex/tiny-bank/internal/domain"

type CreateRequest struct{}

type ReadRequest struct {
	ID string
}

type UpdateBalanceRequest struct {
	ID      string
	Balance int
}

type UpdateTransactionsRequest struct {
	ID          string
	Transaction domain.Transaction
}
