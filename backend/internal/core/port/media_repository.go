package port

import (
	"context"
	"nostalgia/internal/core/domain"
	"time"
)

type MediaRepository interface {
	AddTagToMedia(ctx context.Context, mediaId string, tag string) error
	RemoveTagFromMedia(ctx context.Context, mediaId string, tag string) error
	AddMedia(ctx context.Context, fileId string, thumbnailId string, title string, description string, uploadedBy string, tags []string, uploadedAt time.Time) (string, error)
	GetTagsForMedia(ctx context.Context, mediaId string) ([]string, error)
	GetAllTags(ctx context.Context) ([]string, error)
	GetById(ctx context.Context, id string) (domain.Media, error)
	Get(ctx context.Context) ([]domain.Media, error)
}
