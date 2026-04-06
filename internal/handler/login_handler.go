package handler

import (
	"net/http"

	"github.com/fernando8franco/attengo/internal/apperr"
	"github.com/fernando8franco/attengo/internal/service"
	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	UserService         service.UserService
	RefreshTokenService service.RefreshTokenService
}

func NewLoginHandler(svcUser service.UserService, svcRefreshToken service.RefreshTokenService) *LoginHandler {
	return &LoginHandler{
		UserService:         svcUser,
		RefreshTokenService: svcRefreshToken,
	}
}

type LoginRequest struct {
	Email    string `form:"email"  binding:"required,email"`
	Password string `form:"password"  binding:"required"`
}

func (h *LoginHandler) IndexLogin(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"login.html",
		gin.H{
			"Title": "Login",
		},
	)
}

func (h *LoginHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.Error(apperr.NewBadRequest(err.Error()))
		return
	}

	tokens, err := h.UserService.AdminLogin(c.Request.Context(), service.LoginAdminInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.SetCookie("access_token", tokens.AccessToken, 3600, "/", "", false, true)
	c.SetCookie("refresh_token", tokens.RefreshToken, 86400*7, "/", "", false, true)

	c.Header("HX-Redirect", "/admin/dashboard")
	c.Status(http.StatusOK)
}

func (h *LoginHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	err = h.RefreshTokenService.RevokeToken(c.Request.Context(), refreshToken)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.Redirect(http.StatusFound, "/login")
}
