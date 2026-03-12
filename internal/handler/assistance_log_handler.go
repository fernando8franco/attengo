package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/fernando8franco/attengo/internal/apperr"
	"github.com/fernando8franco/attengo/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AssistanceLogHandler struct {
	assistanceLogService service.AssistanceLogService
}

func NewAssistanceLogHandler(svc service.AssistanceLogService) *AssistanceLogHandler {
	return &AssistanceLogHandler{assistanceLogService: svc}
}

type TakeAttendanceRequest struct {
	UserID       string `json:"user_id"  binding:"required"`
	UserPassword string `json:"user_password"  binding:"required"`
}

func (h *AssistanceLogHandler) TakeAttendance(c *gin.Context) {
	var req TakeAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]string, len(ve))
			for i, fe := range ve {
				out[i] = fmt.Sprintf("field '%s' %s", fe.Field(), fe.Tag())
			}
			c.Error(apperr.NewBadRequest(strings.Join(out, " ")))
			return
		}
		c.Error(apperr.NewBadRequest(err.Error()))
		return
	}

	attendace, err := h.assistanceLogService.TakeAttendance(c.Request.Context(), service.AssistanceLogInput{
		UserID:       req.UserID,
		UserPassword: req.UserPassword,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, attendace)
}
