package port

import (
	"context"
	"io"
)

type GeneratedThumbnail struct {
	File      io.Reader
	MimeType  string
	Extension string
	Blurhash  string
}

type ThumbnailService interface {
	GenerateFromVideo(ctx context.Context, filePath string) (*GeneratedThumbnail, error)
	GenerateBlurHash(file io.Reader, extension string) (string, error)
}
