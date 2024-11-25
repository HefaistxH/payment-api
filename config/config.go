package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	// "github.com/golang-jwt/jwt/v5"
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
	JwtSigningMethod     string
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

	tokenExpire, err := strconv.Atoi(os.Getenv("JWT_EXPIRE_TIME"))
	if err != nil {
		return fmt.Errorf("Error parsing JWT_EXPIRE_TIME: %v", err)
	}
	refreshExpiredToken, err := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRE_TIME"))
	if err != nil {
		return fmt.Errorf("Error parsing JWT_REFRESH_EXPIRE_TIME: %v", err)
	}

	c.JwtConfig = JwtConfig{
		IssuerName:           os.Getenv("JWT_ISSUER_NAME"),
		JwtSignatureKey:      os.Getenv("JWT_SIGNATURE_KEY"),
		JwtSigningMethod:     os.Getenv("JWT_SIGNING_METHOD"),
		JwtExpireTime:        time.Duration(tokenExpire) * time.Minute,
		JwtRefreshExpireTime: time.Duration(refreshExpiredToken) * time.Hour,
	}

	if c.Host == "" || c.Port == "" || c.User == "" || c.Password == "" || c.Name == "" || c.Driver == "" || c.ApiPort == "" || c.IssuerName == "" || c.JwtSignatureKey == "" || c.JwtSigningMethod == "" {
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
