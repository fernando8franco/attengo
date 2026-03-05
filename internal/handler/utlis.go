package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/fernando8franco/attengo/internal/apperr"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"
)

func respondError(c *gin.Context, err error) {
	var notFound *apperr.NotFoundError
	var sqliteErr sqlite3.Error

	switch {
	case errors.As(err, &notFound):
		c.JSON(http.StatusNotFound, gin.H{"error": notFound.Error()})
	case errors.As(err, &sqliteErr):
		c.JSON(http.StatusConflict, gin.H{"error": sqliteErr.Error()})
	default:
		log.Printf("unexpected error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
