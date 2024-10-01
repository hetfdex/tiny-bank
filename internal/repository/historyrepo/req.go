package historyrepo

import "github.com/hetfdex/tiny-bank/internal/domain"

type CreateRequest struct{}

type ReadRequest struct {
	ID string
}

type UpdateRequest struct {
	ID    string
	Event domain.Event
}
