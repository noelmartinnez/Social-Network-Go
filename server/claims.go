package server

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	Username string `json:"username"`
	Rol      string `json:"rol"`
	jwt.RegisteredClaims
}
