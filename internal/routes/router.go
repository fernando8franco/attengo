package routes

import (
	"github.com/fernando8franco/attengo/internal/config"
	"github.com/fernando8franco/attengo/internal/handler"
	"github.com/fernando8franco/attengo/internal/service"
	"github.com/gin-gonic/gin"
)

func SetupRouter(rhSvc service.RequiredHourService, cfg *config.Config) *gin.Engine {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	requiredHoursHandler := handler.NewRequiredHourHandler(rhSvc)

	v1 := r.Group("/api/v1")
	{
		requiredHours := v1.Group("/required_hours")
		{
			requiredHours.POST("", requiredHoursHandler.CreateRequiredHours)
		}
	}

	return r
}
