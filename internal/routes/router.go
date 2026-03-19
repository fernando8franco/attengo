package routes

import (
	"database/sql"
	"reflect"
	"strings"

	"github.com/fernando8franco/attengo/internal/config"
	"github.com/fernando8franco/attengo/internal/handler"
	"github.com/fernando8franco/attengo/internal/middleware"
	"github.com/fernando8franco/attengo/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func initValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

func SetupRouter(conn *sql.DB, cfg *config.Config) *gin.Engine {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	initValidator()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.ErrorHandler())

	rhSvc := service.NewRequiredHourService(conn)
	requiredHoursHandler := handler.NewRequiredHourHandler(rhSvc)

	pSvc := service.NewPeriodService(conn)
	periodHandler := handler.NewPeriodHandler(pSvc)

	uSvc := service.NewUserService(conn, cfg)
	userHandler := handler.NewUserHandler(uSvc)

	alSvc := service.NewAssistanceLogService(conn)
	assistanceLogHandler := handler.NewAssistanceLogHandler(alSvc)

	v1 := r.Group("/api/v1")
	{
		requiredHours := v1.Group("/required_hours")
		{
			requiredHours.POST("", requiredHoursHandler.CreateRequiredHour)
		}

		periods := v1.Group("/periods")
		{
			periods.POST("", periodHandler.CreatePeriod)
		}

		users := v1.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
		}

		setup := v1.Group("/setup")
		{
			setup.POST("/admin", userHandler.SetUpAdmin)
		}

		attendace := v1.Group("/attendace")
		{
			attendace.POST("", assistanceLogHandler.TakeAttendance)
		}
	}

	if cfg.Env == "development" {
		r.POST("/reset", handler.Reset(conn))
	}

	return r
}
