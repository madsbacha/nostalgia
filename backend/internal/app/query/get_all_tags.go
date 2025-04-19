package query

import (
	"context"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/common/errors"
	"nostalgia/internal/core/port"
)

type GetAllTags struct {
}

type GetAllTagsHandler decorator.QueryHandler[GetAllTags, []string]

type getAllTagsHandler struct {
	mediaRepo port.MediaRepository
}

func NewGetAllTagsHandler(mediaRepo port.MediaRepository, logger *logrus.Entry) GetAllTagsHandler {
	return decorator.ApplyQueryDecorators[GetAllTags, []string](
		getAllTagsHandler{mediaRepo: mediaRepo},
		logger,
	)
}

func (h getAllTagsHandler) Handle(ctx context.Context, query GetAllTags) ([]string, error) {
	tags, err := h.mediaRepo.GetAllTags(ctx)
	if err != nil {
		return nil, errors.NewSlugError(err.Error(), "unable-to-get-all-tags")
	}
	return tags, nil
}
