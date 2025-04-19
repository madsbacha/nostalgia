package port

import (
	"context"
)

type Thumbnail struct {
	Id       string
	FileId   string
	Blurhash string
}

type NewThumbnail struct {
	FileId   string
	Blurhash string
}

type ThumbnailRepository interface {
	GetById(ctx context.Context, id string) (*Thumbnail, error)
	Create(ctx context.Context, thumbnail NewThumbnail) (string, error)
}
