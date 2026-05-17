package handler

import (
	"net/http"

	"posixfy-cloud/backend/auth"
	"posixfy-cloud/backend/config"
	"posixfy-cloud/backend/middleware"
	"posixfy-cloud/backend/models"
	"posixfy-cloud/backend/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	UserService *service.UserService
	Config      *config.Config
}

func NewAuthHandler(us *service.UserService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{UserService: us, Config: cfg}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.UserService.Authenticate(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := auth.GenerateToken(h.Config.JWTSecret, user.ID, user.Username, user.UID, user.GID, user.Groups, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		User:  *user,
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	claims := middleware.GetClaims(c)
	user, err := h.UserService.GetByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
