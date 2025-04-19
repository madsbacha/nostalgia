package query

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/config"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/common/errors"
	token2 "nostalgia/internal/common/token"
	"time"
)

type GetTokenForUser struct {
	UserId string
}

type GetTokenForUserHandler decorator.QueryHandler[GetTokenForUser, string]

type getTokenForUserHandler struct {
	jwtConfig config.Jwt
}

func NewGetTokenForUserHandler(jwtConfig config.Jwt, logger *logrus.Entry) GetTokenForUserHandler {
	return decorator.ApplyQueryDecorators[GetTokenForUser, string](
		getTokenForUserHandler{
			jwtConfig: jwtConfig,
		},
		logger,
	)
}

func (h getTokenForUserHandler) Handle(ctx context.Context, q GetTokenForUser) (string, error) {
	claims := token2.UserTokenClaims{
		UserId: q.UserId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(h.jwtConfig.HmacSecret)

	if err != nil {
		return "", errors.NewSlugError(err.Error(), "unable-to-sign-token")
	}

	return tokenString, nil
}
