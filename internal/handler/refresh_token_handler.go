package handler

import (
	"net/http"

	"github.com/fernando8franco/attengo/internal/apperr"
	"github.com/fernando8franco/attengo/internal/auth"
	"github.com/fernando8franco/attengo/internal/service"
	"github.com/gin-gonic/gin"
)

type RefreshTokenHandler struct {
	RefreshTokenService service.RefreshTokenService
}

func NewRefreshTokenHandler(svc service.RefreshTokenService) *RefreshTokenHandler {
	return &RefreshTokenHandler{RefreshTokenService: svc}
}

func (h *RefreshTokenHandler) RefreshAccessToken(c *gin.Context) {
	refreshToken, err := auth.GetBearerToken(c.Request.Header)
	if err != nil {
		c.Error(apperr.NewBadRequest("Couldn't find the refresh token"))
		return
	}

	accessToken, err := h.RefreshTokenService.CreateAccessToken(c.Request.Context(), refreshToken)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, accessToken)
}
