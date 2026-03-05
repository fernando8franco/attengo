package handler

import (
	"net/http"

	"github.com/fernando8franco/attengo/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{UserService: svc}
}

type CreateUserRequest struct {
	Name           string `json:"name"  binding:"required"`
	Email          string `json:"email"  binding:"required,email"`
	RequiredHourID int    `json:"required_hour_id"  binding:"required"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": err.Error()})
		return
	}

	u, err := h.UserService.CreateUser(c.Request.Context(), service.CreateUserInput{
		Name:           req.Email,
		Email:          req.Email,
		RequiredHourID: req.RequiredHourID,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, u)
}
