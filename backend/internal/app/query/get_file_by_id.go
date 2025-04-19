package query

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/config"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/common/errors"
	token2 "nostalgia/internal/common/token"
	"nostalgia/internal/core/domain"
	"nostalgia/internal/core/port"
	"time"
)

type GetFileById struct {
	Id string
}

type GetFileByIdHandler decorator.QueryHandler[GetFileById, *domain.File]

type getFileByIdHandler struct {
	fileRepo  port.FileRepository
	jwtConfig config.Jwt
}

func NewGetFileById(fileRepo port.FileRepository, jwtConfig config.Jwt, logger *logrus.Entry) GetFileByIdHandler {
	return decorator.ApplyQueryDecorators[GetFileById, *domain.File](
		getFileByIdHandler{fileRepo: fileRepo, jwtConfig: jwtConfig},
		logger,
	)
}

func (h getFileByIdHandler) Handle(ctx context.Context, q GetFileById) (*domain.File, error) {
	file, err := h.fileRepo.GetById(ctx, q.Id)
	if err != nil {
		return nil, err
	}

	claims := token2.FileTokenClaims{
		FileId: file.Id,
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
		return nil, errors.NewSlugError(err.Error(), "unable-to-sign-token")
	}

	publicFilePath := fmt.Sprintf("/api/files/%s?jwt=%s", file.Id, tokenString)

	return &domain.File{
		Id:           file.Id,
		MimeType:     file.MimeType,
		Path:         publicFilePath,
		InternalPath: file.Path,
	}, nil
}
