package handler

import (
	"errors"
	"net/http"

	"github.com/fernando8franco/attengo/internal/service"
	"github.com/gin-gonic/gin"

	"github.com/mattn/go-sqlite3"
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

	rh, err := h.requiredHourService.CreateRequiredHour(c.Request.Context(), service.CrateRequiredHourInput{
		Type:    req.Type,
		Minutes: req.Minutes,
	})

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				c.JSON(http.StatusBadRequest, gin.H{"error": "This record already exists"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": sqliteErr.ExtendedCode})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rh)
}
