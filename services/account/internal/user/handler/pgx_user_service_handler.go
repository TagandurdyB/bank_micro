package handler

import (
	"bank_micro/services/account/internal/user/model"
	"bank_micro/services/account/internal/user/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PgxUserHandler struct {
	service *service.PgxUserService
}

func NewPgxUserHandler(svc *service.PgxUserService) *PgxUserHandler {
	return &PgxUserHandler{service: svc}
}

func (h *PgxUserHandler) RegisterUserRoutes(r *gin.Engine) {
	r.GET("/api/user", h.GetAllUsers)                     //✅
	r.GET("/api/user/:id", h.GetUserByID)                 //✅
	r.POST("/api/user", h.CreateUser)                     //✅
	r.PUT("/api/user/:id", h.UpdateUser)                  //✅
	r.DELETE("/api/user/:id", h.SoftDeleteUser)           //✅
	r.POST("/api/user/:id/deposit", h.Deposit)            //✅
	r.PATCH("/api/user/:id/reference", h.UpdateReference) //✅
}

func (h *PgxUserHandler) GetAllUsers(c *gin.Context) {
	loadAccount, _ := strconv.ParseBool(c.DefaultQuery("load_account", "false"))

	users, err := h.service.GetAllUsers(loadAccount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]gin.H, 0, len(users))
	for _, u := range users {
		resp = append(resp, h.mapUser(&u))
	}

	c.JSON(http.StatusOK, resp)
}

func (h *PgxUserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	loadAccount, _ := strconv.ParseBool(c.DefaultQuery("load_account", "false"))

	user, err := h.service.GetUserByID(id, loadAccount)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.mapUser(user))
}

func (h *PgxUserHandler) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := h.service.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, h.mapUser(createdUser))
}

func (h *PgxUserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = id

	updatedUser, err := h.service.UpdateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.mapUser(updatedUser))
}

func (h *PgxUserHandler) SoftDeleteUser(c *gin.Context) {
	id := c.Param("id")

	deletedUser, err := h.service.SoftDeleteUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.mapUser(deletedUser))
}

func (h *PgxUserHandler) Deposit(c *gin.Context) {
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

	err := h.service.DepositToUser(id, amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deposit successful", "user_id": id})
}

func (h *PgxUserHandler) UpdateReference(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id is required"})
		return
	}

	var body map[string]any
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	refUserID, ok := body["ref_user_id"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ref_user_id is required or invalid"})
		return
	}

	err := h.service.UpdateUserReference(userID, &refUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user reference updated successfully",
		"user_id": userID,
		"new_ref": refUserID,
	})
}

func (h *PgxUserHandler) mapUser(u *model.User) gin.H {
	userMap := gin.H{
		"id":          u.ID,
		"account_id":  u.AccountID,
		"ref_user_id": u.RefUserID,
		"is_locked":   u.IsLocked,
		"created_at":  u.CreatedAt,
		"deleted_at":  u.DeletedAt,
	}

	if u.Account != nil {
		userMap["account"] = gin.H{
			"id":         u.Account.ID,
			"balance":    u.Account.Balance,
			"currency":   u.Account.Currency,
			"is_locked":  u.Account.IsLocked,
			"created_at": u.Account.CreatedAt,
		}
	}

	return userMap
}
