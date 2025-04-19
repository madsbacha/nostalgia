package token

import "github.com/golang-jwt/jwt/v5"

type FileTokenClaims struct {
	FileId string `json:"file_id"`
	jwt.RegisteredClaims
}
