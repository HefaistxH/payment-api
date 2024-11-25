package service

import (
	"fmt"
	"mnc-techtest/config"
	"mnc-techtest/entity"
	"mnc-techtest/entity/dto"
	"mnc-techtest/shared/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	GenerateToken(credentials entity.Credential) (dto.AuthResponseDto, error)
	GenerateRefreshToken(credentials entity.Credential) (dto.AuthResponseDto, error)
	ParseToken(tokenHeader string) (jwt.MapClaims, error)
}

type jwtService struct {
	cfg config.JwtConfig
}

func (j *jwtService) GenerateToken(credentials entity.Credential) (dto.AuthResponseDto, error) {
	accessClaims := model.MyCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.cfg.IssuerName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.JwtExpireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		CredId: credentials.Id,
		Email:  credentials.Email,
		Role:   credentials.Role,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(j.cfg.JwtSignatureKey)

	if err != nil {
		return dto.AuthResponseDto{}, fmt.Errorf("failed to generate token: %v", err)
	}

	return dto.AuthResponseDto{
		Token: accessTokenString,
	}, nil
}

func (j *jwtService) GenerateRefreshToken(credentials entity.Credential) (dto.AuthResponseDto, error) {
	refreshClaims := model.MyCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.cfg.IssuerName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.JwtRefreshExpireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Email: credentials.Email,
		Role:  credentials.Role,
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(j.cfg.JwtSignatureKey)

	if err != nil {
		return dto.AuthResponseDto{}, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	return dto.AuthResponseDto{
		RefreshToken: refreshTokenString,
	}, nil
}

func (j *jwtService) ParseToken(tokenHeader string) (jwt.MapClaims, error) {
	tokenParsed, err := jwt.ParseWithClaims(tokenHeader, &model.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.cfg.JwtSignatureKey), nil
	})

	if err != nil {
		return nil, err
	}

	claim, ok := tokenParsed.Claims.(*model.MyCustomClaims)
	if !ok || !tokenParsed.Valid {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	return jwt.MapClaims{
		"email": claim.Email,
		"role":  claim.Role,
	}, nil
}

func NewJwtService(cfg config.JwtConfig) JwtService {
	return &jwtService{cfg: cfg}
}
