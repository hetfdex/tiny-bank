package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hetfdex/tiny-bank/internal/service"
)

const (
	baseURL = "/api/v1/users/"
)

type Handler interface {
	ConfigHandlers(router *gin.Engine)
}

type hdl struct {
	svc service.Service
}

func New(svc service.Service) Handler {
	return &hdl{
		svc: svc,
	}
}

func (h hdl) ConfigHandlers(router *gin.Engine) {
	router.POST(baseURL, h.createUser)
	router.POST(baseURL+":user_id", h.createAccount)
	router.DELETE(baseURL+":user_id", h.deactivateUser)
	router.PUT(baseURL+":user_id/accounts/:account_id", h.deposit)
	router.PATCH(baseURL+":user_id/accounts/:account_id", h.withdraw)
	router.POST(baseURL+":user_id/accounts/:account_id", h.transfer)
	router.GET(baseURL+":user_id/accounts/:account_id", h.balance)
	router.GET(baseURL+":user_id/accounts/:account_id/transactions", h.transactions)
}

func (h hdl) createUser(c *gin.Context) {
	var req service.CreateUserRequest

	err := c.BindJSON(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	res, err := h.svc.CreateUser(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h hdl) createAccount(c *gin.Context) {
	res, err := h.svc.CreateAccount(
		service.CreateAccountRequest{
			UserID: c.Param("user_id"),
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h hdl) deactivateUser(c *gin.Context) {
	err := h.svc.DeactivateUser(
		service.DeactivateUserRequest{
			UserID: c.Param("user_id"),
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h hdl) deposit(c *gin.Context) {
	req := service.DepositRequest{}

	err := c.BindJSON(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	req.UserID = c.Param("user_id")
	req.AccountID = c.Param("account_id")

	res, err := h.svc.Deposit(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (h hdl) withdraw(c *gin.Context) {
	req := service.WithdrawRequest{}

	err := c.BindJSON(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	req.UserID = c.Param("user_id")
	req.AccountID = c.Param("account_id")

	res, err := h.svc.Withdraw(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (h hdl) transfer(c *gin.Context) {
	req := service.TransferRequest{}

	err := c.BindJSON(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	req.SenderUserID = c.Param("user_id")
	req.SenderAccountID = c.Param("account_id")

	res, err := h.svc.Transfer(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (h hdl) balance(c *gin.Context) {
	res, err := h.svc.Balance(
		service.BalanceRequest{
			UserID:    c.Param("user_id"),
			AccountID: c.Param("account_id"),
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (h hdl) transactions(c *gin.Context) {
	res, err := h.svc.Transactions(
		service.TransactionsRequest{
			UserID:    c.Param("user_id"),
			AccountID: c.Param("account_id"),
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}
