package query

import (
	"context"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/core/domain"
	"nostalgia/internal/core/port"
)

type GetDefaultThumbnail struct {
}

type GetDefaultThumbnailHandler decorator.QueryHandler[GetDefaultThumbnail, domain.Thumbnail]

type getDefaultThumbnailHandler struct {
	settingRepo   port.SettingRepository
	thumbnailRepo port.ThumbnailRepository
}

func NewGetDefaultThumbnailHandler(settingRepo port.SettingRepository, thumbnailRepo port.ThumbnailRepository, logger *logrus.Entry) GetDefaultThumbnailHandler {
	return decorator.ApplyQueryDecorators[GetDefaultThumbnail, domain.Thumbnail](
		getDefaultThumbnailHandler{settingRepo: settingRepo, thumbnailRepo: thumbnailRepo},
		logger,
	)
}

func (h getDefaultThumbnailHandler) Handle(ctx context.Context, query GetDefaultThumbnail) (domain.Thumbnail, error) {
	defaultThumbnailSetting, err := h.settingRepo.GetByKey(ctx, domain.SettingKeyDefaultThumbnail)
	if err != nil {
		return domain.Thumbnail{}, err
	}

	thumbnail, err := h.thumbnailRepo.GetById(ctx, defaultThumbnailSetting.Value)
	if err != nil {
		return domain.Thumbnail{}, err
	}

	return domain.Thumbnail{
		Id:       thumbnail.Id,
		FileId:   thumbnail.FileId,
		BlurHash: thumbnail.Blurhash,
	}, nil
}
