package handler

import (
	"net/http"

	"github.com/fernando8franco/attengo/internal/apperr"
	"github.com/fernando8franco/attengo/internal/auth"
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

	url, err := h.UserService.CreateUser(c.Request.Context(), service.CreateUserInput{
		Name:           req.Name,
		Email:          req.Email,
		RequiredHourID: req.RequiredHourID,
		PeriodID:       req.PeriodID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("Location", url)

	c.Status(http.StatusCreated)
}

type SetUpAdminRequest struct {
	Name     string `json:"name"  binding:"required"`
	Email    string `json:"email"  binding:"required,email"`
	Password string `json:"password"  binding:"required"`
}

func (h *UserHandler) SetUpAdmin(c *gin.Context) {
	var req SetUpAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperr.NewBadRequest(err.Error()))
		return
	}

	admin, err := h.UserService.SetUpAdmin(c.Request.Context(), service.CreateAdminInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, admin)
}

type LoginRequest struct {
	Email    string `json:"email"  binding:"required,email"`
	Password string `json:"password"  binding:"required"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperr.NewBadRequest(err.Error()))
		return
	}

	tokens, err := h.UserService.AdminLogin(c.Request.Context(), service.LoginAdminInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *UserHandler) Logout(c *gin.Context) {
	refreshToken, err := auth.GetBearerToken(c.Request.Header)
	if err != nil {
		c.Error(apperr.NewBadRequest("Couldn't find the refresh token"))
		return
	}

	err = h.UserService.AdminLogout(c.Request.Context(), refreshToken)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
