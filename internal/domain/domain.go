package domain

import "time"

type User struct {
	ID         string
	CreatedAt  time.Time
	Active     bool
	Name       string
	AccountIDs map[string]struct{}
}

type Account struct {
	ID           string
	CreatedAt    time.Time
	Balance      int
	Transactions []Transaction
}

type Transaction struct {
	Timestamp         time.Time
	Operation         string
	Amount            int
	ReceiverUserID    string
	SenderUserID      string
	ReceiverAccountID string
	SenderAccountID   string
}
