package handler

import (
	"bank_micro/services/account/internal/account/model"
	"bank_micro/services/account/internal/account/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PgxAccountHandler struct {
	service *service.PgxAccountService
}

func NewPgxAccountHandler(svc *service.PgxAccountService) *PgxAccountHandler {
	return &PgxAccountHandler{service: svc}
}

func (h *PgxAccountHandler) RegisterAccountRoutes(r *gin.Engine) {
	r.GET("/api/account", h.GetAllAccounts)
	r.GET("/api/account/:id", h.GetAccountByID)
	r.POST("/api/account", h.CreateAccount)
	r.PUT("/api/account/:id", h.UpdateAccount)
	r.DELETE("/api/account/:id", h.SoftDeleteAccount)
	r.PATCH("/api/account/:id/balance", h.UpdateBalance)
}

func (h *PgxAccountHandler) GetAllAccounts(c *gin.Context) {
	balanceStr := c.Query("balance")
	var balance *int64
	if balanceStr != "" {
		val, err := strconv.ParseInt(balanceStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid balance"})
			return
		}
		balance = &val
	}

	accounts, err := h.service.GetAllAccounts(balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]gin.H, 0, len(accounts))
	for _, acc := range accounts {
		resp = append(resp, h.mapAccount(&acc))
	}

	c.JSON(http.StatusOK, resp)
}

func (h *PgxAccountHandler) GetAccountByID(c *gin.Context) {
	id := c.Param("id")

	acc, err := h.service.GetAccountByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.mapAccount(acc))
}

func (h *PgxAccountHandler) CreateAccount(c *gin.Context) {
	var acc model.Account
	if err := c.ShouldBindJSON(&acc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdAcc, err := h.service.CreateAccount(&acc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, h.mapAccount(createdAcc))
}

func (h *PgxAccountHandler) UpdateAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}
	var acc model.Account
	if err := c.ShouldBindJSON(&acc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	acc.ID = id

	updatedAcc, err := h.service.UpdateAccount(&acc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.mapAccount(updatedAcc))
}

func (h *PgxAccountHandler) SoftDeleteAccount(c *gin.Context) {
	id := c.Param("id")

	deletedAcc, err := h.service.SoftDeleteAccount(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.mapAccount(deletedAcc))
}

func (h *PgxAccountHandler) UpdateBalance(c *gin.Context) {
	id := c.Param("id")

	var body map[string]any
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	amountFloat, ok := body["amount"].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount is required or invalid"})
		return
	}
	amount := int64(amountFloat)

	newBalance, err := h.service.UpdateAccountBalance(id, amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "balance": newBalance})
}

func (h *PgxAccountHandler) mapAccount(acc *model.Account) gin.H {
	return gin.H{
		"id":         acc.ID,
		"balance":    acc.Balance,
		"currency":   acc.Currency,
		"is_locked":  acc.IsLocked,
		"created_at": acc.CreatedAt,
		"deleted_at": acc.DeletedAt,
	}
}
