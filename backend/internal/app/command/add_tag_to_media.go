package command

import (
	"context"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/common/errors"
	"nostalgia/internal/core/port"
)

type AddTagToMedia struct {
	Tag     string
	MediaId string
}

type AddTagToMediaHandler decorator.CommandHandler[AddTagToMedia]

type addTagToMediaHandler struct {
	mediaRepo port.MediaRepository
}

func NewAddTagToMediaHandler(mediaRepo port.MediaRepository, logger *logrus.Entry) AddTagToMediaHandler {
	return decorator.ApplyCommandDecorators[AddTagToMedia](
		addTagToMediaHandler{mediaRepo: mediaRepo},
		logger,
	)
}

func (h addTagToMediaHandler) Handle(ctx context.Context, cmd AddTagToMedia) error {
	err := h.mediaRepo.AddTagToMedia(ctx, cmd.Tag, cmd.MediaId)
	if err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-add-tag-to-media")
	}
	return nil
}
