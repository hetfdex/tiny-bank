package userrepo

type CreateRequest struct {
	Name string
}

type ReadRequest struct {
	ID string
}

type UpdateStatusRequest struct {
	ID     string
	Active bool
}

type UpdateAccountIDsRequest struct {
	ID        string
	AccountID string
}
