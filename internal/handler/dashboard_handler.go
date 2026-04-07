package handler

import (
	"net/http"

	"github.com/fernando8franco/attengo/internal/service"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	AssistanceLogService service.AssistanceLogService
}

func NewDashboardHandler(svc service.AssistanceLogService) *DashboardHandler {
	return &DashboardHandler{AssistanceLogService: svc}
}

func (h *DashboardHandler) Index(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"dashboard.html",
		gin.H{
			"Title": "Dashboard",
		},
	)
}
