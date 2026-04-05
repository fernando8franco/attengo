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

	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("web/templates/**/*.html")

	rhSvc := service.NewRequiredHourService(conn)
	requiredHoursHandler := handler.NewRequiredHourHandler(rhSvc)

	pSvc := service.NewPeriodService(conn)
	periodHandler := handler.NewPeriodHandler(pSvc)

	uSvc := service.NewUserService(conn, cfg)
	userHandler := handler.NewUserHandler(uSvc)

	alSvc := service.NewAssistanceLogService(conn)
	assistanceLogHandler := handler.NewAssistanceLogHandler(alSvc)

	rtSvc := service.NewRefreshTokenService(conn, cfg)
	refreshTokenHandler := handler.NewRefreshTokenHandler(rtSvc)

	r.GET("/", assistanceLogHandler.Index)
	r.POST("/attendace", assistanceLogHandler.Attendance)

	v1 := r.Group("/api/v1")
	{
		requiredHours := v1.Group("/required_hours")
		requiredHours.Use(middleware.AuthRequired(cfg.IssuerJWT, cfg.SecretJWT))
		{
			requiredHours.POST("", requiredHoursHandler.CreateRequiredHour)
		}

		periods := v1.Group("/periods")
		periods.Use(middleware.AuthRequired(cfg.IssuerJWT, cfg.SecretJWT))
		{
			periods.POST("", periodHandler.CreatePeriod)
		}

		users := v1.Group("/users")
		users.Use(middleware.AuthRequired(cfg.IssuerJWT, cfg.SecretJWT))
		{
			users.POST("", userHandler.CreateUser)
		}

		setup := v1.Group("/setup")
		{
			setup.POST("/admin", userHandler.SetUpAdmin)
		}

		/* attendace := v1.Group("/attendace")
		{
			attendace.POST("", assistanceLogHandler.TakeAttendance)
		} */

		v1.POST("/login", userHandler.Login)
		v1.POST("/logout", userHandler.Logout)
		v1.POST("/refresh", refreshTokenHandler.RefreshAccessToken)
	}

	if cfg.Env == "development" {
		r.POST("/reset", handler.Reset(conn))
	}

	return r
}
