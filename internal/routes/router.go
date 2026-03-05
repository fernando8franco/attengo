package routes

import (
	"database/sql"

	"github.com/fernando8franco/attengo/internal/config"
	"github.com/fernando8franco/attengo/internal/handler"
	"github.com/fernando8franco/attengo/internal/service"
	"github.com/gin-gonic/gin"
)

func SetupRouter(conn *sql.DB, cfg *config.Config) *gin.Engine {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	rhSvc := service.NewRequiredHourService(conn)
	uSvc := service.NewUserService(conn)
	requiredHoursHandler := handler.NewRequiredHourHandler(rhSvc)
	userHandler := handler.NewUserHandler(uSvc)

	v1 := r.Group("/api/v1")
	{
		requiredHours := v1.Group("/required_hours")
		{
			requiredHours.POST("", requiredHoursHandler.CreateRequiredHours)
		}

		users := v1.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
		}
	}

	if cfg.Env == "development" {
		r.POST("/reset", handler.Reset(conn))
	}

	return r
}
