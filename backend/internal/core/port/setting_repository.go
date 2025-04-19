package port

import (
	"context"
	"nostalgia/internal/core/domain"
)

type SettingRepository interface {
	GetByKey(ctx context.Context, key string) (domain.Setting, error)
	Set(ctx context.Context, key, value string) error
}
