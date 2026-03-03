package handler

import (
	"net/http"

	"github.com/fernando8franco/attengo/internal/repository"
	"github.com/gin-gonic/gin"
)

type requiredHoursHandler struct {
	queries repository.Queries
}

func NewRequiredHoursHandler(q repository.Queries) *requiredHoursHandler {
	return &requiredHoursHandler{queries: q}
}

type CreateRequiredHoursRequest struct {
	Type    string `json:"type"  binding:"required"`
	Minutes int    `json:"minutes"  binding:"required"`
}

func (h *requiredHoursHandler) CreateRequiredHours(c *gin.Context) {
	var req CreateRequiredHoursRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": err.Error()})
		return
	}

	rh, err := h.queries.CreateRequiredHours(c.Request.Context(), repository.CreateRequiredHoursParams{
		Type:    req.Type,
		Minutes: int64(req.Minutes),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create required hours"})
		return
	}

	c.JSON(http.StatusCreated, rh)
}
