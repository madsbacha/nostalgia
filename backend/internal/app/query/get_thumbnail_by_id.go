package query

import (
	"context"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/core/domain"
	"nostalgia/internal/core/port"
)

type GetThumbnailById struct {
	Id string
}

type GetThumbnailByIdHandler decorator.QueryHandler[GetThumbnailById, *domain.Thumbnail]

type getThumbnailByIdHandler struct {
	thumbnailRepo port.ThumbnailRepository
}

func NewGetThumbnailById(thumbnailRepo port.ThumbnailRepository, logger *logrus.Entry) GetThumbnailByIdHandler {
	return decorator.ApplyQueryDecorators[GetThumbnailById, *domain.Thumbnail](
		getThumbnailByIdHandler{thumbnailRepo: thumbnailRepo},
		logger,
	)
}

func (h getThumbnailByIdHandler) Handle(ctx context.Context, q GetThumbnailById) (*domain.Thumbnail, error) {
	thumbnail, err := h.thumbnailRepo.GetById(ctx, q.Id)
	if err != nil {
		return nil, err
	}

	return &domain.Thumbnail{
		Id:       thumbnail.Id,
		FileId:   thumbnail.FileId,
		BlurHash: thumbnail.Blurhash,
	}, nil
}
