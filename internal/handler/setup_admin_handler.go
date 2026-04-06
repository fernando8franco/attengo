package handler

import (
	"net/http"

	"github.com/fernando8franco/attengo/internal/apperr"
	"github.com/fernando8franco/attengo/internal/service"
	"github.com/gin-gonic/gin"
)

type SetUpAdminHandler struct {
	UserService service.UserService
}

func NewSetUpAdminHandler(svc service.UserService) *SetUpAdminHandler {
	return &SetUpAdminHandler{UserService: svc}
}

func (h *SetUpAdminHandler) IndexSetUpAdmin(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"setup-admin.html",
		gin.H{
			"Title": "SetUp Admin",
		},
	)
}

type SetUpAdminRequest struct {
	Name     string `form:"name"  binding:"required"`
	Email    string `form:"email"  binding:"required,email"`
	Password string `form:"password"  binding:"required"`
}

func (h *SetUpAdminHandler) SetUpAdmin(c *gin.Context) {
	var req SetUpAdminRequest
	if err := c.ShouldBind(&req); err != nil {
		c.Error(apperr.NewBadRequest(err.Error()))
		return
	}

	tokens, err := h.UserService.SetUpAdmin(c.Request.Context(), service.CreateAdminInput{
		Name:     req.Name,
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
