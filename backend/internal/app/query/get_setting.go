package query

import (
	"context"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/core/domain"
	"nostalgia/internal/core/port"
)

type GetSetting struct {
	Key string
}

type GetSettingHandler decorator.QueryHandler[GetSetting, domain.Setting]

type getSettingHandler struct {
	settingRepo port.SettingRepository
}

func NewGetSetting(settingRepo port.SettingRepository, logger *logrus.Entry) GetSettingHandler {
	return decorator.ApplyQueryDecorators[GetSetting, domain.Setting](
		getSettingHandler{settingRepo: settingRepo},
		logger,
	)
}

func (h getSettingHandler) Handle(ctx context.Context, q GetSetting) (domain.Setting, error) {
	return h.settingRepo.GetByKey(ctx, q.Key)
}
