package middleware

import (
	"net/http"

	"github.com/fernando8franco/attengo/internal/auth"
	"github.com/gin-gonic/gin"
)

type contextKey string

const userIDKey contextKey = "userID"

func AuthRequired(issuer, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := auth.GetBearerToken(c.Request.Header)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		userID, err := auth.ValidateJWT(issuer, secret, accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Set(string(userIDKey), userID)
		c.Next()
	}
}

func GetUserID(c *gin.Context) (string, bool) {
	val, exists := c.Get(string(userIDKey))
	if !exists {
		return "", false
	}
	userID, ok := val.(string)
	return userID, ok
}

func AuthMiddleware(issuer, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("access_token")
		if err != nil {
			if c.GetHeader("HX-Request") == "true" {
				c.Header("HX-Redirect", "/login")
				c.AbortWithStatus(http.StatusUnauthorized)
			} else {
				c.Redirect(http.StatusFound, "/login")
				c.Abort()
			}
			return
		}

		claims, err := auth.ValidateJWT(issuer, secret, token)
		if err != nil {
			c.SetCookie("access_token", "", -1, "/", "", false, true)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Set("userID", claims)
		c.Next()
	}
}
