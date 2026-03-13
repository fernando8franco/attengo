package handler

import (
	"net/http"

	"github.com/fernando8franco/attengo/internal/apperr"
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
	RequiredHourID int    `json:"required_hour_id"  binding:"required,numeric"`
	PeriodID       int    `json:"period_id"  binding:"required,numeric"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperr.NewBadRequest(err.Error()))
		return
	}

	u, err := h.UserService.CreateUser(c.Request.Context(), service.CreateUserInput{
		Name:           req.Name,
		Email:          req.Email,
		RequiredHourID: req.RequiredHourID,
		PeriodID:       req.PeriodID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, u)
}
