package handler

import (
	"database/sql"

	"github.com/fernando8franco/attengo/internal/repository"
	"github.com/gin-gonic/gin"
)

func Reset(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		q := repository.New(db)
		err := q.DeleteRequiredHours(c.Request.Context())
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to reset database"})
			return
		}
		c.JSON(200, gin.H{"message": "Database reset successfully"})
	}
}
