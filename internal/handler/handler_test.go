package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hetfdex/tiny-bank/internal/domain"
	"github.com/hetfdex/tiny-bank/internal/service"
	"github.com/hetfdex/tiny-bank/test/mock/servicemock"
	"github.com/stretchr/testify/assert"
)

func TestCreate_ErrJSON(t *testing.T) {
	httpReq := makeHTTPRequest(
		t,
		http.MethodPost,
		baseURL,
		nil,
	)

	hdl := New(&servicemock.Mock{})

	rr := setupTest(hdl, httpReq)

	assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	assert.Equal(t, "{\"error\":\"invalid request\"}", rr.Body.String())
}

func TestCreateUser_ErrCreate(t *testing.T) {
	req := service.CreateUserRequest{
		Name: "joe",
	}

	httpReq := makeHTTPRequest(
		t,
		http.MethodPost,
		baseURL,
		makeBody(req),
	)

	svc := &servicemock.Mock{}

	svc.On(
		"CreateUser",
		req,
	).Return(
		service.CreateUserResponse{},
		errors.New("error create"),
	)

	hdl := New(svc)

	rr := setupTest(hdl, httpReq)

	assert.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode)
	assert.Equal(t, "{\"error\":\"error create\"}", rr.Body.String())
}

func TestCreateUser_Ok(t *testing.T) {
	req := service.CreateUserRequest{
		Name: "joe",
	}

	httpReq := makeHTTPRequest(
		t,
		http.MethodPost,
		baseURL,
		makeBody(req),
	)

	svc := &servicemock.Mock{}

	svc.On(
		"CreateUser",
		req,
	).Return(
		service.CreateUserResponse{
			UserID: "1",
		},
		nil,
	)

	hdl := New(svc)

	rr := setupTest(hdl, httpReq)

	assert.Equal(t, http.StatusCreated, rr.Result().StatusCode)
	assert.Equal(t, "{\"user_id\":\"1\"}", rr.Body.String())
}

func TestCreateAccount_ErrCreate(t *testing.T) {
	httpReq := makeHTTPRequest(
		t,
		http.MethodPost,
		baseURL+"1",
		nil,
	)

	svc := &servicemock.Mock{}

	svc.On(
		"CreateAccount",
		service.CreateAccountRequest{
			UserID: "1",
		},
	).Return(
		service.CreateAccountResponse{},
		errors.New("error create"),
	)

	hdl := New(svc)

	rr := setupTest(hdl, httpReq)

	assert.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode)
	assert.Equal(t, "{\"error\":\"error create\"}", rr.Body.String())
}

func TestCreateAccount_Ok(t *testing.T) {
	httpReq := makeHTTPRequest(
		t,
		http.MethodPost,
		baseURL+"1",
		nil,
	)

	svc := &servicemock.Mock{}

	svc.On(
		"CreateAccount",
		service.CreateAccountRequest{
			UserID: "1",
		},
	).Return(
		service.CreateAccountResponse{
			AccountID: "2",
		},
		nil,
	)

	hdl := New(svc)

	rr := setupTest(hdl, httpReq)

	assert.Equal(t, http.StatusCreated, rr.Result().StatusCode)
	assert.Equal(t, "{\"account_id\":\"2\"}", rr.Body.String())
}

func TestDeactivateUser_ErrDeactivate(t *testing.T) {
	req := service.DeactivateUserRequest{
		UserID: "1",
	}

	httpReq := makeHTTPRequest(
		t,
		http.MethodDelete,
		baseURL+"1",
		nil,
	)

	svc := &servicemock.Mock{}

	svc.On(
		"DeactivateUser",
		req,
	).Return(
		errors.New("error deactivate"),
	)

	hdl := New(svc)

	rr := setupTest(hdl, httpReq)

	assert.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode)
	assert.Equal(t, "{\"error\":\"error deactivate\"}", rr.Body.String())
}

func TestDeactivateUser_Ok(t *testing.T) {
	req := service.DeactivateUserRequest{
		UserID: "1",
	}

	httpReq := makeHTTPRequest(
		t,
		http.MethodDelete,
		baseURL+"1",
		nil,
	)

	svc := &servicemock.Mock{}

	svc.On(
		"DeactivateUser",
		req,
	).Return(
		nil,
	)

	hdl := New(svc)

	rr := setupTest(hdl, httpReq)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	assert.Equal(t, "{\"status\":\"ok\"}", rr.Body.String())
}

func TestDeposit_ErrJSON(t *testing.T) {
	httpReq := makeHTTPRequest(
		t,
		http.MethodPut,
		baseURL+"1/accounts/2",
		nil,
	)

	hdl := New(&servicemock.Mock{})

	rr := setupTest(hdl, httpReq)

	assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	assert.Equal(t, "{\"error\":\"invalid request\"}", rr.Body.String())
}

func TestDeposit_ErrDeposit(t *testing.T) {
	httpReq := makeHTTPRequest(
		t,
		http.MethodPut,
		baseURL+"1/accounts/2",
		makeBody(
			service.DepositRequest{
				Amount: 3,
			},
		),
	)

	svc := &servicemock.Mock{}

	svc.On(
		"Deposit",
		service.DepositRequest{
			UserID:    "1",
			AccountID: "2",
			Amount:    3,
		},
	).Return(
		service.DepositResponse{},
		errors.New("error deposit"),
	)

	hdl := New(svc)

	rr := setupTest(hdl, httpReq)

	assert.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode)
	assert.Equal(t, "{\"error\":\"error deposit\"}", rr.Body.String())
}

func TestDeposit_Ok(t *testing.T) {
	httpReq := makeHTTPRequest(
		t,
		http.MethodPut,
		baseURL+"1/accounts/2",
		makeBody(
			service.DepositRequest{
				Amount: 3,
			},
		),
	)

	svc := &servicemock.Mock{}

	svc.On(
		"Deposit",
		service.DepositRequest{
			UserID:    "1",
			AccountID: "2",
			Amount:    3,
		},
	).Return(
		service.DepositResponse{
			Balance: 10,
		},
		nil,
	)

	hdl := New(svc)

	rr := setupTest(hdl, httpReq)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	assert.Equal(t, "{\"balance\":10}", rr.Body.String())
}

func TestWithdraw_ErrJSON(t *testing.T) {
	httpReq := makeHTTPRequest(
		t,
		http.MethodPatch,
		baseURL+"1/accounts/2",
		nil,
	)

	hdl := New(&servicemock.Mock{})

	rr := setupTest(hdl, httpReq)

	assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	assert.Equal(t, "{\"error\":\"invalid request\"}", rr.Body.String())
}

func TestWithdraw_ErrWithdraw(t *testing.T) {
	httpReq := makeHTTPRequest(
		t,
		http.MethodPatch,
		baseURL+"1/accounts/2",
		makeBody(
			service.WithdrawRequest{
				Amount: 3,
			},
		),
	)

	svc := &servicemock.Mock{}

	svc.On(
		"Withdraw",
		service.WithdrawRequest{
			UserID:    "1",
			AccountID: "2",
			Amount:    3,
		},
	).Return(
		service.WithdrawResponse{},
		errors.New("error withdraw"),
	)

	hdl := New(svc)

	rr := setupTest(hdl, httpReq)

	assert.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode)
	assert.Equal(t, "{\"error\":\"error withdraw\"}", rr.Body.String())
}

func TestWithdraw_Ok(t *testing.T) {
	httpReq := makeHTTPRequest(
		t,
		http.MethodPatch,
		baseURL+"1/accounts/2",
		makeBody(
			service.WithdrawRequest{
				Amount: 3,
			},
		),
	)

	svc := &servicemock.Mock{}

	svc.On(
		"Withdraw",
		service.WithdrawRequest{
			UserID:    "1",
			AccountID: "2",
			Amount:    3,
		},
	).Return(
		service.WithdrawResponse{
			Balance: 10,
		},
		nil,
	)

	hdl := New(svc)

	rr := setupTest(hdl, httpReq)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	assert.Equal(t, "{\"balance\":10}", rr.Body.String())
}

func TestTransfer_ErrJSON(t *testing.T) {
	httpReq := makeHTTPRequest(
		t,
		http.MethodPost,
		baseURL+"1/accounts/2",
		nil,
	)

	hdl := New(&servicemock.Mock{})

	rr := setupTest(hdl, httpReq)

	assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	assert.Equal(t, "{\"error\":\"invalid request\"}", rr.Body.String())
}

func TestTransfer_ErrTransfer(t *testing.T) {
	httpReq := makeHTTPRequest(
		t,
		http.MethodPost,
		baseURL+"1/accounts/3",
		makeBody(
			service.TransferRequest{
				ReceiverUserID:    "2",
				ReceiverAccountID: "4",
				Amount:            5,
			},
		),
	)

	svc := &servicemock.Mock{}

	svc.On(
		"Transfer",
		service.TransferRequest{
			SenderUserID:      "1",
			ReceiverUserID:    "2",
			SenderAccountID:   "3",
			ReceiverAccountID: "4",
			Amount:            5,
		},
	).Return(
		service.TransferResponse{},
		errors.New("error transfer"),
	)

	hdl := New(svc)

	rr := setupTest(hdl, httpReq)

	assert.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode)
	assert.Equal(t, "{\"error\":\"error transfer\"}", rr.Body.String())
}

func TestTransfer_Ok(t *testing.T) {
	httpReq := makeHTTPRequest(
		t,
		http.MethodPost,
		baseURL+"1/accounts/3",
		makeBody(
			service.TransferRequest{
				ReceiverUserID:    "2",
				ReceiverAccountID: "4",
				Amount:            5,
			},
		),
	)

	svc := &servicemock.Mock{}

	svc.On(
		"Transfer",
		service.TransferRequest{
			SenderUserID:      "1",
			ReceiverUserID:    "2",
			SenderAccountID:   "3",
			ReceiverAccountID: "4",
			Amount:            5,
		},
	).Return(
		service.TransferResponse{
			Balance: 10,
		},
		nil,
	)

	hdl := New(svc)

	rr := setupTest(hdl, httpReq)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	assert.Equal(t, "{\"balance\":10}", rr.Body.String())
}

func TestBalance_ErrBalance(t *testing.T) {
	httpReq := makeHTTPRequest(
		t,
		http.MethodGet,
		baseURL+"1/accounts/2",
		nil,
	)

	svc := &servicemock.Mock{}

	svc.On(
		"Balance",
		service.BalanceRequest{
			UserID:    "1",
			AccountID: "2",
		},
	).Return(
		service.BalanceResponse{},
		errors.New("error balance"),
	)

	hdl := New(svc)

	rr := setupTest(hdl, httpReq)

	assert.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode)
	assert.Equal(t, "{\"error\":\"error balance\"}", rr.Body.String())
}

func TestBalance_Ok(t *testing.T) {
	httpReq := makeHTTPRequest(
		t,
		http.MethodGet,
		baseURL+"1/accounts/2",
		nil,
	)

	svc := &servicemock.Mock{}

	svc.On(
		"Balance",
		service.BalanceRequest{
			UserID:    "1",
			AccountID: "2",
		},
	).Return(
		service.BalanceResponse{
			Balance: 10,
		},
		nil,
	)

	hdl := New(svc)

	rr := setupTest(hdl, httpReq)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	assert.Equal(t, "{\"balance\":10}", rr.Body.String())
}

func TestTransactions_ErrTransaction(t *testing.T) {
	httpReq := makeHTTPRequest(
		t,
		http.MethodGet,
		baseURL+"1/accounts/2/transactions",
		nil,
	)

	svc := &servicemock.Mock{}

	svc.On(
		"Transactions",
		service.TransactionsRequest{
			UserID:    "1",
			AccountID: "2",
		},
	).Return(
		service.TransactionsResponse{},
		errors.New("error transaction"),
	)

	hdl := New(svc)

	rr := setupTest(hdl, httpReq)

	assert.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode)
	assert.Equal(t, "{\"error\":\"error transaction\"}", rr.Body.String())
}

func TestTransactions_Ok(t *testing.T) {
	httpReq := makeHTTPRequest(
		t,
		http.MethodGet,
		baseURL+"1/accounts/2/transactions",
		nil,
	)

	svc := &servicemock.Mock{}

	svc.On(
		"Transactions",
		service.TransactionsRequest{
			UserID:    "1",
			AccountID: "2",
		},
	).Return(
		service.TransactionsResponse{
			Transactions: []domain.Transaction{
				{
					Timestamp:         time.Time{},
					Operation:         "operation",
					Amount:            666,
					ReceiverUserID:    "1",
					SenderUserID:      "2",
					ReceiverAccountID: "3",
					SenderAccountID:   "4",
				},
			},
		},
		nil,
	)

	hdl := New(svc)

	rr := setupTest(hdl, httpReq)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	assert.Equal(t, "{\"transactions\":[{\"Timestamp\":\"0001-01-01T00:00:00Z\",\"Operation\":\"operation\",\"Amount\":666,\"ReceiverUserID\":\"1\",\"SenderUserID\":\"2\",\"ReceiverAccountID\":\"3\",\"SenderAccountID\":\"4\"}]}", rr.Body.String())
}

func setupTest(hdl Handler, req *http.Request) *httptest.ResponseRecorder {
	router := gin.Default()

	hdl.ConfigHandlers(router)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	return rr
}

func makeHTTPRequest(
	t *testing.T,
	httpMethod string,
	url string,
	body io.Reader,
) *http.Request {
	req, err := http.NewRequest(
		httpMethod,
		url,
		body,
	)

	if err != nil {
		t.Fatal(err)
	}

	return req
}

func makeBody(v interface{}) io.Reader {
	body, err := json.Marshal(v)

	if err != nil {
		return nil
	}

	return bytes.NewBuffer(body)
}
