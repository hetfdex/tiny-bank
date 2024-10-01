package service

import "github.com/hetfdex/tiny-bank/internal/domain"

type CreateResponse struct {
	ID         string
	AccountIDs []string
}

type BalanceResponse struct {
	Balance int `json:"balance"`
}

type DepositResponse BalanceResponse

type WithdrawResponse DepositResponse

type TransferResponse DepositResponse

type HistoryResponse struct {
	Events []domain.Event `json:"events"`
}
