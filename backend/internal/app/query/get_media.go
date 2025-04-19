package query

import (
	"context"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/core/domain"
	"nostalgia/internal/core/port"
)

type GetMedia struct {
}

type GetMediaHandler decorator.QueryHandler[GetMedia, []domain.Media]

type getMediaHandler struct {
	mediaRepo port.MediaRepository
}

func NewGetMedia(mediaRepo port.MediaRepository, logger *logrus.Entry) GetMediaHandler {
	return decorator.ApplyQueryDecorators[GetMedia, []domain.Media](
		getMediaHandler{mediaRepo: mediaRepo},
		logger,
	)
}

func (h getMediaHandler) Handle(ctx context.Context, q GetMedia) ([]domain.Media, error) {
	mediaList, err := h.mediaRepo.Get(ctx)
	if err != nil {
		return nil, err
	}

	return mediaList, nil
}
