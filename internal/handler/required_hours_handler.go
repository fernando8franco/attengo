package handler

import (
	"net/http"

	"github.com/fernando8franco/attengo/internal/service"
	"github.com/gin-gonic/gin"
)

type RequiredHoursHandler struct {
	requiredHourService service.RequiredHourService
}

func NewRequiredHourHandler(svc service.RequiredHourService) *RequiredHoursHandler {
	return &RequiredHoursHandler{requiredHourService: svc}
}

type CreateRequiredHoursRequest struct {
	Type    string `json:"type"  binding:"required"`
	Minutes int    `json:"minutes"  binding:"required"`
}

func (h *RequiredHoursHandler) CreateRequiredHours(c *gin.Context) {
	var req CreateRequiredHoursRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": err.Error()})
		return
	}

	rh, err := h.requiredHourService.CreateRequiredHour(c.Request.Context(), service.CreateRequiredHourInput{
		Type:         req.Type,
		TotalMinutes: req.Minutes,
	})

	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, rh)
}
