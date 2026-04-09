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

func (h *UserHandler) Index(c *gin.Context) {
	hours, periods, users, err := h.UserService.GetHoursPeriodsAndUsers(c.Request.Context())
	if err != nil {
		c.Error(apperr.NewBadRequest(err.Error()))
		return
	}

	c.HTML(
		http.StatusOK,
		"users.html",
		gin.H{
			"Title":     "Usuarios",
			"HoursType": hours,
			"Periods":   periods,
			"Users":     users,
		},
	)
}

type CreateUserRequest struct {
	Name           string `form:"name"  binding:"required"`
	Email          string `form:"email"  binding:"required,email"`
	RequiredHourID int    `form:"required_hour_id"  binding:"required,numeric"`
	PeriodID       int    `form:"period_id"  binding:"required,numeric"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBind(&req); err != nil {
		c.Error(apperr.NewBadRequest(err.Error()))
		return
	}

	user, err := h.UserService.CreateUser(c.Request.Context(), service.CreateUserInput{
		Name:           req.Name,
		Email:          req.Email,
		RequiredHourID: req.RequiredHourID,
		PeriodID:       req.PeriodID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.HTML(
		http.StatusOK,
		"users-info-row",
		user,
	)
}

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
