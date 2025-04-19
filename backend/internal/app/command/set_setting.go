package command

import (
	"context"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/core/port"
)

type SetSetting struct {
	Key   string
	Value string
}

type SetSettingHandler decorator.CommandHandler[SetSetting]

type setSettingHandler struct {
	settingRepo port.SettingRepository
}

func NewSetSettingHandler(settingRepo port.SettingRepository, logger *logrus.Entry) SetSettingHandler {
	return decorator.ApplyCommandDecorators[SetSetting](
		setSettingHandler{settingRepo: settingRepo},
		logger,
	)
}

func (h setSettingHandler) Handle(ctx context.Context, cmd SetSetting) error {
	return h.settingRepo.Set(ctx, cmd.Key, cmd.Value)
}
