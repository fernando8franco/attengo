package handler

import (
	"bytes"
	"html/template"
	"net/http"
	"time"

	"github.com/fernando8franco/attengo/internal/apperr"
	"github.com/fernando8franco/attengo/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService service.UserService
	Templates   *template.Template
}

func NewUserHandler(svc service.UserService, tmpl *template.Template) *UserHandler {
	return &UserHandler{
		UserService: svc,
		Templates:   tmpl,
	}
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

/* type SetUpAdminRequest struct {
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
} */

// type LoginRequest struct {
// 	Email    string `json:"email"  binding:"required,email"`
// 	Password string `json:"password"  binding:"required"`
// }

// func (h *UserHandler) Login(c *gin.Context) {
// 	var req LoginRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.Error(apperr.NewBadRequest(err.Error()))
// 		return
// 	}

// 	tokens, err := h.UserService.AdminLogin(c.Request.Context(), service.LoginAdminInput{
// 		Email:    req.Email,
// 		Password: req.Password,
// 	})
// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, tokens)
// }

// func (h *UserHandler) Logout(c *gin.Context) {
// 	refreshToken, err := auth.GetBearerToken(c.Request.Header)
// 	if err != nil {
// 		c.Error(apperr.NewBadRequest("Couldn't find the refresh token"))
// 		return
// 	}

// 	err = h.UserService.AdminLogout(c.Request.Context(), refreshToken)
// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	c.Status(http.StatusNoContent)
// }

func (h *UserHandler) StreamUserHandler(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	clientGone := c.Request.Context().Done()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	sendUsersUpdate := func() bool {
		users, err := h.UserService.GetActiveUsers(c.Request.Context())
		if err != nil {
			return false
		}

		var buf bytes.Buffer
		err = h.Templates.ExecuteTemplate(&buf, "view-users", gin.H{"Users": users})
		if err != nil {
			return false
		}

		c.SSEvent("message", buf.String())
		c.Writer.Flush()
		return true
	}

	if ok := sendUsersUpdate(); !ok {
		return
	}

	for {
		select {
		case <-clientGone:
			return
		case <-ticker.C:
			if ok := sendUsersUpdate(); !ok {
				return
			}
		}
	}
}
