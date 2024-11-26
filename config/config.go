package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Driver   string
}

type ApiConfig struct {
	ApiPort string
}

type JwtConfig struct {
	IssuerName           string
	JwtSignatureKey      string
	JwtSigningMethod     *jwt.SigningMethodHMAC
	JwtExpireTime        time.Duration
	JwtRefreshExpireTime time.Duration
}

type Config struct {
	DbConfig
	ApiConfig
	JwtConfig
}

func (c *Config) readConfig() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf(("Error loading .env file: %v"), err)
	}

	c.DbConfig = DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		Driver:   os.Getenv("DB_DRIVER"),
	}
	c.ApiConfig = ApiConfig{ApiPort: os.Getenv("API_PORT")}

	tokenExpire, err := strconv.Atoi(os.Getenv("TOKEN_EXPIRE"))
	if err != nil {
		return fmt.Errorf("Error parsing JWT_EXPIRE_TIME: %v", err)
	}
	refreshExpiredToken, err := strconv.Atoi(os.Getenv("TOKEN_REFRESH_EXPIRE"))
	if err != nil {
		return fmt.Errorf("Error parsing JWT_REFRESH_EXPIRE_TIME: %v", err)
	}

	c.JwtConfig = JwtConfig{
		IssuerName:           os.Getenv("TOKEN_ISSUE"),
		JwtSignatureKey:      os.Getenv("TOKEN_SECRET"),
		JwtSigningMethod:     jwt.SigningMethodHS256,
		JwtExpireTime:        time.Duration(tokenExpire) * time.Minute,
		JwtRefreshExpireTime: time.Duration(refreshExpiredToken) * time.Hour,
	}

	if c.Host == "" || c.Port == "" || c.User == "" || c.Password == "" || c.Name == "" || c.Driver == "" || c.ApiPort == "" || c.IssuerName == "" || c.JwtSignatureKey == "" {
		return fmt.Errorf("Missing required environment variables for configuration")
	}

	return nil
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := cfg.readConfig(); err != nil {
		return nil, err
	}
	return cfg, nil
}
