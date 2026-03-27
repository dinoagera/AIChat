package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/dinoagera/AIChat/internal/domain"
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
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			au.log.Info("failed to register user", "err", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, Response{Message: messages.MsgInvalidCredentials})
			return
		}
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
	accessToken, refreshToken, err := au.authService.SignIn(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		//TODO: Add process if wrong password or email not exist; Можно не возвращать, т.к. злоумышленник не поймет, что подбирает пароли под существующий акк
		au.log.Info("failed to login user", "err", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, ResponseWithTokens{AccessToken: accessToken, RefreshToken: refreshToken})
}
func (au *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		au.log.Info("failed to decode json req", "err", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	accessToken, refreshToken, err := au.authService.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		au.log.Info("failed to refresh tokens", "err", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, ResponseWithTokens{AccessToken: accessToken, RefreshToken: refreshToken})
}
func (h *AuthHandler) SetupRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.SignUp)
		auth.POST("/login", h.SignIn)
		auth.POST("/refresh", h.Refresh)
	}
}
