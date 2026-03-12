package handler

import (
	"net/http"

	"github.com/fernando8franco/attengo/internal/apperr"
	"github.com/fernando8franco/attengo/internal/service"
	"github.com/gin-gonic/gin"
)

type PeriodHandler struct {
	periodService service.PeriodService
}

func NewPeriodHandler(svc service.PeriodService) *PeriodHandler {
	return &PeriodHandler{periodService: svc}
}

type CreatePeriodRequest struct {
	Name      string `json:"name"  binding:"required"`
	EntryDate string `json:"entry_date" binding:"required,datetime=2006-01-02"`
	ExitDate  string `json:"exit_date"  binding:"required,datetime=2006-01-02"`
}

func (h *PeriodHandler) CreatePeriod(c *gin.Context) {
	var req CreatePeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperr.NewBadRequest(err.Error()))
		return
	}

	period, err := h.periodService.CreatePeriod(c.Request.Context(), service.CreatePeriodInput{
		Name:      req.Name,
		EntryDate: req.EntryDate,
		ExitDate:  req.ExitDate,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, period)
}
