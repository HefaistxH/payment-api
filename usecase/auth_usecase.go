package usecase

import (
	"fmt"
	"mnc-techtest/entity/dto"
	"mnc-techtest/repository"
	"mnc-techtest/shared/service"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Login(payload dto.AuthRequestDto) (dto.AuthResponseDto, error)
	RefreshToken(refreshToken string) (dto.AuthResponseDto, error)
	Logout(token string) error
}

type authUsecase struct {
	customerRepo repository.CustomerRepository
	credRepo     repository.CredentialRepository
	jwtService   service.JwtService
}

func (a *authUsecase) Login(payload dto.AuthRequestDto) (dto.AuthResponseDto, error) {
	// Check user existence
	isExist, err := a.customerRepo.CheckCustomerByEmail(payload.Email)
	if err != nil || !isExist {
		return dto.AuthResponseDto{}, fmt.Errorf("user not found or invalid credentials")
	}
	// Fetch user and compare password
	user, err := a.credRepo.GetCredByEmail(payload.Email)
	if err != nil {
		return dto.AuthResponseDto{}, err
	}
	logrus.Infof("User found: %v", user)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return dto.AuthResponseDto{}, fmt.Errorf("invalid credentials")
	}
	// Generate JWT token
	accessToken, err := a.jwtService.GenerateToken(user)
	if err != nil {
		logrus.Warnf("Failed to generate access token: %v", err)
		return dto.AuthResponseDto{}, nil
	}
	refreshToken, err := a.jwtService.GenerateRefreshToken(user)
	if err != nil {
		logrus.Warnf("Failed to generate refresh token: %v", err)
		return dto.AuthResponseDto{}, nil
	}
	authResponse := dto.AuthResponseDto{
		Token:        accessToken.Token,
		RefreshToken: refreshToken.RefreshToken,
	}
	return authResponse, nil
}

func (a *authUsecase) RefreshToken(refreshToken string) (dto.AuthResponseDto, error) {
	claims, err := a.jwtService.ParseToken(refreshToken)
	if err != nil {
		logrus.Warnf("Failed to parse refresh token: %v", err)
		return dto.AuthResponseDto{}, fmt.Errorf("failed to parse refresh token: %v", err)
	}
	email, ok := claims["email"].(string)
	if !ok || email == "" {
		logrus.Warn("Email claim is missing or invalid")
		return dto.AuthResponseDto{}, fmt.Errorf("email claim is missing or invalid")
	}
	user, err := a.credRepo.GetCredByEmail(email)
	if err != nil {
		return dto.AuthResponseDto{}, err
	}
	newToken, err := a.jwtService.GenerateToken(user)
	if err != nil {
		logrus.Warnf("Failed to generate new token: %v", err)
		return dto.AuthResponseDto{}, err
	}
	logrus.Infof("Token refreshed: %v", newToken)
	return dto.AuthResponseDto{
		Token:        newToken.Token,
		RefreshToken: refreshToken,
	}, nil
}

func (a *authUsecase) Logout(token string) error {
	err := a.jwtService.InvalidateToken(token)
	if err != nil {
		return fmt.Errorf("Failed to invalidate token: %v", err)
	}
	return nil
}

func NewAuthUsecase(customerRepo repository.CustomerRepository, credRepo repository.CredentialRepository, jwtService service.JwtService) AuthUsecase {
	return &authUsecase{
		customerRepo: customerRepo,
		credRepo:     credRepo,
		jwtService:   jwtService}
}
