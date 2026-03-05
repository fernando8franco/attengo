package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/fernando8franco/attengo/internal/apperr"
	"github.com/gin-gonic/gin"
)

func respondError(c *gin.Context, err error) {
	var notFound *apperr.NotFoundError
	var conflict *apperr.ConflictError

	switch {
	case errors.As(err, &notFound):
		c.JSON(http.StatusNotFound, gin.H{"error": notFound.Error()})
	case errors.As(err, &conflict):
		c.JSON(http.StatusConflict, gin.H{"error": conflict.Error()})
	default:
		log.Printf("unexpected error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
