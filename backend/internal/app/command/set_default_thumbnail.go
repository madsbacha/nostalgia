package command

import (
	"context"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/core/domain"
	"nostalgia/internal/core/port"
)

type SetDefaultThumbnail struct {
	ThumbnailId string
}

type SetDefaultThumbnailHandler decorator.CommandHandler[SetDefaultThumbnail]

type setDefaultThumbnailHandler struct {
	settingRepo port.SettingRepository
}

func NewSetDefaultThumbnailHandler(settingRepo port.SettingRepository, logger *logrus.Entry) SetDefaultThumbnailHandler {
	return decorator.ApplyCommandDecorators[SetDefaultThumbnail](
		setDefaultThumbnailHandler{settingRepo: settingRepo},
		logger,
	)
}

func (h setDefaultThumbnailHandler) Handle(ctx context.Context, cmd SetDefaultThumbnail) error {
	return h.settingRepo.Set(ctx, domain.SettingKeyDefaultThumbnail, cmd.ThumbnailId)
}
