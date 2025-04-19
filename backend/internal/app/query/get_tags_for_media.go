package query

import (
	"context"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/common/errors"
	"nostalgia/internal/core/port"
)

type GetTagsForMedia struct {
	MediaId string
}

type GetTagsForMediaHandler decorator.QueryHandler[GetTagsForMedia, []string]

type getTagsForMediaHandler struct {
	mediaRepo port.MediaRepository
}

func NewGetTagsForMediaHandler(mediaRepo port.MediaRepository, logger *logrus.Entry) GetTagsForMediaHandler {
	return decorator.ApplyQueryDecorators[GetTagsForMedia, []string](
		getTagsForMediaHandler{mediaRepo: mediaRepo},
		logger,
	)
}

func (h getTagsForMediaHandler) Handle(ctx context.Context, query GetTagsForMedia) ([]string, error) {
	tags, err := h.mediaRepo.GetTagsForMedia(ctx, query.MediaId)
	if err != nil {
		return nil, errors.NewSlugError(err.Error(), "unable-to-get-tags-for-media")
	}
	return tags, nil
}
