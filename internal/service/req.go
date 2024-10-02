package service

type CreateUserRequest struct {
	Name string `json:"name"`
}

type CreateAccountRequest struct {
	UserID string `json:"user_id"`
}

type DeactivateUserRequest struct {
	UserID string `json:"user_id"`
}

type DepositRequest struct {
	UserID    string `json:"user_id"`
	AccountID string `json:"account_id"`
	Amount    int    `json:"amount"`
}

type WithdrawRequest struct {
	UserID    string `json:"user_id"`
	AccountID string `json:"account_id"`
	Amount    int    `json:"amount"`
}

type TransferRequest struct {
	SenderUserID      string `json:"sender_user_id"`
	ReceiverUserID    string `json:"receiver_user_id"`
	SenderAccountID   string `json:"sender_account_id"`
	ReceiverAccountID string `json:"receiver_account_id"`
	Amount            int    `json:"amount"`
}

type BalanceRequest struct {
	UserID    string `json:"user_id"`
	AccountID string `json:"account_id"`
}

type TransactionsRequest struct {
	UserID    string `json:"user_id"`
	AccountID string `json:"account_id"`
}
