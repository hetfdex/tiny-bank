package accountrepo

type CreateRequest struct {
	HistoryID string
}

type ReadRequest struct {
	ID string
}

type UpdateRequest struct {
	ID      string
	Balance int
}
