package handler

import (
	"net/http"

	"github.com/fernando8franco/attengo/internal/apperr"
	"github.com/fernando8franco/attengo/internal/service"
	"github.com/gin-gonic/gin"
)

type RequiredHourHandler struct {
	requiredHourService service.RequiredHourService
}

func NewRequiredHourHandler(svc service.RequiredHourService) *RequiredHourHandler {
	return &RequiredHourHandler{requiredHourService: svc}
}

type CreateRequiredHourRequest struct {
	Type    string `json:"type"  binding:"required"`
	Minutes int    `json:"minutes"  binding:"required"`
}

func (h *RequiredHourHandler) CreateRequiredHour(c *gin.Context) {
	var req CreateRequiredHourRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperr.NewBadRequest(err.Error()))
		return
	}

	rh, err := h.requiredHourService.CreateRequiredHour(c.Request.Context(), service.RequiredHourInput{
		Type:         req.Type,
		TotalMinutes: req.Minutes,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, rh)
}
