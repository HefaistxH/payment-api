package service

import (
	"fmt"
	"mnc-techtest/config"
	"mnc-techtest/entity"
	"mnc-techtest/entity/dto"
	"mnc-techtest/shared/model"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type JwtService interface {
	GenerateToken(credentials entity.Credential) (dto.AuthResponseDto, error)
	GenerateRefreshToken(credentials entity.Credential) (dto.AuthResponseDto, error)
	ParseToken(tokenHeader string) (jwt.MapClaims, error)
	InvalidateToken(tokenHeader string) error
}

type jwtService struct {
	cfg            config.JwtConfig
	invalidateList map[string]time.Time
	mu             sync.Mutex
}

func NewJwtService(cfg config.JwtConfig) JwtService {
	return &jwtService{
		cfg:            cfg,
		invalidateList: make(map[string]time.Time),
	}
}

func (j *jwtService) GenerateToken(credentials entity.Credential) (dto.AuthResponseDto, error) {
	accessClaims := model.MyCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.cfg.IssuerName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.JwtExpireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		CredId: credentials.Id,
		UserID: credentials.UserId,
		Email:  credentials.Email,
		Role:   credentials.Role,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(j.cfg.JwtSignatureKey))

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
		CredId: credentials.Id,
		UserID: credentials.UserId,
		Email:  credentials.Email,
		Role:   credentials.Role,
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(j.cfg.JwtSignatureKey))

	if err != nil {
		return dto.AuthResponseDto{}, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	return dto.AuthResponseDto{
		RefreshToken: refreshTokenString,
	}, nil
}

func (j *jwtService) ParseToken(tokenHeader string) (jwt.MapClaims, error) {
	if invalidateTime, found := j.invalidateList[tokenHeader]; found && time.Now().Before(invalidateTime) {
		logrus.Infof("Token is invalidated: %s", tokenHeader)
		return nil, fmt.Errorf("token is invalidated")
	}

	tokenParsed, err := jwt.ParseWithClaims(tokenHeader, &model.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.cfg.JwtSignatureKey), nil
	})

	if err != nil {
		logrus.Errorf("Error parsing token: %v", err)
		return nil, err
	}

	if claims, ok := tokenParsed.Claims.(*model.MyCustomClaims); ok && tokenParsed.Valid {
		logrus.Infof("Token claims: %+v", claims)
		return jwt.MapClaims{
			"credId":  claims.CredId,
			"user_id": claims.UserID, // Ensure that the claim is UserID
			"email":   claims.Email,
			"role":    claims.Role,
		}, nil
	} else {
		logrus.Errorf("Failed to parse token: %v", err)
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}
}

func (j *jwtService) InvalidateToken(tokenHeader string) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.invalidateList[tokenHeader] = time.Now().Add(j.cfg.JwtExpireTime)
	return nil
}
