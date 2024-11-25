package controller

import (
	"mnc-techtest/config"
	"mnc-techtest/entity/dto"
	"mnc-techtest/usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUsecase usecase.AuthUsecase
	rg          *gin.RouterGroup
}

func NewAuthController(authUsecase usecase.AuthUsecase, rg *gin.RouterGroup) AuthController {
	return AuthController{authUsecase: authUsecase, rg: rg}
}

func (a *AuthController) Login(ctx *gin.Context) {
	var payload dto.AuthRequestDto
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	authResponse, err := a.authUsecase.Login(payload)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, authResponse)
}

func (a *AuthController) RefreshToken(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing Authorization header"})
		return
	}

	// Extract token from header
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Authorization header format"})
		return
	}

	authResponse, err := a.authUsecase.RefreshToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, authResponse)
}

func (a *AuthController) Logout(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing Authorization header"})
		return
	}
	// Extract token from header
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Authorization header format"})
		return
	}
	err := a.authUsecase.Logout(token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (a *AuthController) Route() {
	a.rg.POST(config.LoginRoute, a.Login)
	a.rg.POST(config.RefreshRoute, a.RefreshToken)
	a.rg.POST(config.LogoutRoute, a.Logout)
}
