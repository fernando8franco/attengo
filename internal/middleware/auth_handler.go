package middleware

import (
	"net/http"

	"github.com/fernando8franco/attengo/internal/auth"
	"github.com/fernando8franco/attengo/internal/service"
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

func AuthMiddleware(issuer, secret string, refreshTokenService service.RefreshTokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("access_token")

		if err == nil {
			claims, err := auth.ValidateJWT(issuer, secret, token)
			if err == nil {
				c.Set("userID", claims)
				c.Next()
				return
			}
		}

		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			redirectUnauthorized(c)
			return
		}

		tokenInfo, err := refreshTokenService.CreateAccessToken(c.Request.Context(), refreshToken)
		if err != nil {
			redirectUnauthorized(c)
			return
		}

		c.SetCookie("access_token", tokenInfo.AccessToken, 3600, "/", "", false, true)
		c.Set("userID", tokenInfo.UserID)
		c.Next()
	}
}

func redirectUnauthorized(c *gin.Context) {
	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Redirect", "/login")
		c.AbortWithStatus(http.StatusUnauthorized)
	} else {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
	}
}
