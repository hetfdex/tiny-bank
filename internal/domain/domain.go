package domain

import "time"

type User struct {
	ID         string
	CreatedAt  time.Time
	Active     bool
	Name       string
	AccountIDs []string
}

type Account struct {
	ID        string
	CreatedAt time.Time
	Balance   int
	HistoryID string
}

type History struct {
	ID        string
	CreatedAt time.Time
	Events    []Event
}

type Event struct {
	Timestamp         time.Time
	Operation         string
	Amount            int
	ReceiverUserID    string
	SenderUserID      string
	ReceiverAccountID string
	SenderAccountID   string
}
