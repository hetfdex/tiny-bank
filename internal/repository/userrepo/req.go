package userrepo

type CreateRequest struct {
	Name      string
	AccountID string
}

type ReadRequest struct {
	ID string
}

type UpdateRequest struct {
	ID     string
	Active bool
}
