package handler

import (
	"log/slog"
	"net/http"

	"github.com/dinoagera/AIChat/pkg/messages"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	log         *slog.Logger
	authService AuthService
}

func NewAuthHandler(log *slog.Logger, authService AuthService) *AuthHandler {
	return &AuthHandler{
		log:         log,
		authService: authService,
	}
}
func (au *AuthHandler) SignUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		au.log.Info("failed to decode json req", "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, Response{Message: messages.MsgInvalidCredentials})
		return
	}
	err := au.authService.SignUp(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		au.log.Info("failed to register user", "err", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, Response{Message: messages.MsgUserCreated})
}
func (au *AuthHandler) SignIn(c *gin.Context) {
	var req SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		au.log.Info("failed to decode json req", "err", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	accessjwt, refreshToken, err := au.authService.SignIn(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		//TODO: Add process to wrong password or email not exist
		au.log.Info("failed to login user", "err", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, ResponseSignUp{AccessToken: accessjwt, RefreshToken: refreshToken})
}
func (h *AuthHandler) SetupRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.SignUp)
		auth.POST("/login", h.SignIn)
	}
}
