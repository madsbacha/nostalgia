package query

import (
	"context"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/core/domain"
	"nostalgia/internal/core/port"
)

type GetMediaById struct {
	Id string
}

type GetMediaByIdHandler decorator.QueryHandler[GetMediaById, domain.Media]

type getMediaByIdHandler struct {
	mediaRepo port.MediaRepository
}

func NewGetMediaById(mediaRepo port.MediaRepository, logger *logrus.Entry) GetMediaByIdHandler {
	return decorator.ApplyQueryDecorators[GetMediaById, domain.Media](
		getMediaByIdHandler{mediaRepo: mediaRepo},
		logger,
	)
}

func (h getMediaByIdHandler) Handle(ctx context.Context, q GetMediaById) (domain.Media, error) {
	media, err := h.mediaRepo.GetById(ctx, q.Id)
	if err != nil {
		return domain.Media{}, err
	}

	return media, nil
}
