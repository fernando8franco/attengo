package middleware

import (
	"net/http"

	"github.com/fernando8franco/attengo/internal/apperr"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		if appErr, ok := err.(*apperr.ErrorResponse); ok {
			c.JSON(appErr.Code, appErr)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "An unexpected error occurred",
		})
	}
}
