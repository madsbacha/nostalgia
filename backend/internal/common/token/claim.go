package token

import "github.com/golang-jwt/jwt/v5"

type UserTokenClaims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}
