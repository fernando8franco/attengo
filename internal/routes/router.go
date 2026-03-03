package routes

import (
	"github.com/fernando8franco/attengo/internal/config"
	"github.com/fernando8franco/attengo/internal/handler"
	"github.com/fernando8franco/attengo/internal/repository"
	"github.com/gin-gonic/gin"
)

func SetupRouter(q repository.Queries, cfg *config.Config) *gin.Engine {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	requiredHoursHandler := handler.NewRequiredHoursHandler(q)

	v1 := r.Group("/api/v1")
	{
		requiredHours := v1.Group("/required_hours")
		{
			requiredHours.POST("", requiredHoursHandler.CreateRequiredHours)
		}
	}

	return r
}
