package handler

import (
	"log/slog"

	"github.com/dinoagera/AIChat/internal/http/response"
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
		c.AbortWithStatus(400)
		return
	}
	err := au.authService.SignUp(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		au.log.Info("failed to register user", "err", err)
		c.AbortWithStatusJSON(400, response.Error(err.Error()))
		return
	}
	c.Status(201)
}
func (au *AuthHandler) SignIn(c *gin.Context) {
	var req SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		au.log.Info("failed to decode json req", "err", err)
		c.AbortWithStatus(400)
		return
	}
	err := au.authService.SignUp(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		au.log.Info("failed to register user", "err", err)
		c.AbortWithStatusJSON(400, response.Error(err.Error()))
		return
	}
	c.Status(201)
}
