package middleware

import (
	"github.com/fernando8franco/attengo/internal/auth"
	"github.com/gin-gonic/gin"
)

func AuthRequired(issuer, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := auth.GetBearerToken(c.Request.Header)
		if err != nil {
			c.Error(err)
			return
		}

		userID, err := auth.ValidateJWT(accessToken, issuer, secret)
		if err != nil {
			c.Error(err)
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
