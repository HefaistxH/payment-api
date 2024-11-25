package model

import "github.com/golang-jwt/jwt/v5"

type MyCustomClaims struct {
	jwt.RegisteredClaims
	CredId string `json:"credId"`
	Email  string `json:"userId"`
	Role   string `json:"role"`
}
