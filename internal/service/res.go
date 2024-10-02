package service

import "github.com/hetfdex/tiny-bank/internal/domain"

type CreateUserResponse struct {
	UserID string `json:"user_id"`
}

type CreateAccountResponse struct {
	AccountID string `json:"account_id"`
}

type BalanceResponse struct {
	Balance int `json:"balance"`
}

type DepositResponse BalanceResponse

type WithdrawResponse DepositResponse

type TransferResponse DepositResponse

type TransactionsResponse struct {
	Transactions []domain.Transaction `json:"transactions"`
}
