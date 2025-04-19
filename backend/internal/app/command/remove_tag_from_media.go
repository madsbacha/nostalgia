package command

import (
	"context"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/common/errors"
	"nostalgia/internal/core/port"
)

type RemoveTagFromMedia struct {
	Tag     string
	MediaId string
}

type RemoveTagFromMediaHandler decorator.CommandHandler[RemoveTagFromMedia]

type removeTagFromMediaHandler struct {
	mediaRepo port.MediaRepository
}

func NewRemoveTagFromMediaHandler(mediaRepo port.MediaRepository, logger *logrus.Entry) RemoveTagFromMediaHandler {
	return decorator.ApplyCommandDecorators[RemoveTagFromMedia](
		removeTagFromMediaHandler{mediaRepo: mediaRepo},
		logger,
	)
}

func (h removeTagFromMediaHandler) Handle(ctx context.Context, cmd RemoveTagFromMedia) error {
	err := h.mediaRepo.RemoveTagFromMedia(ctx, cmd.Tag, cmd.MediaId)
	if err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-remove-tag-from-media")
	}
	return nil
}
