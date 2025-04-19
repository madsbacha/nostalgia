package request

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/core/port"
)

type AddThumbnail struct {
	File      io.Reader
	Extension string
}

type AddThumbnailHandler decorator.RequestHandler[AddThumbnail, string]

type addThumbnailHandler struct {
	fileRepo        port.FileRepository
	thumbnailRepo   port.ThumbnailRepository
	thumbnailServce port.ThumbnailService
}

func NewAddThumbnailHandler(fileRepo port.FileRepository, thumbnailRepo port.ThumbnailRepository, thumbnailService port.ThumbnailService, logger *logrus.Entry) AddThumbnailHandler {
	return decorator.ApplyRequestDecorators[AddThumbnail, string](
		addThumbnailHandler{fileRepo: fileRepo, thumbnailRepo: thumbnailRepo, thumbnailServce: thumbnailService},
		logger,
	)
}

func (h addThumbnailHandler) Handle(ctx context.Context, cmd AddThumbnail) (string, error) {
	mimeType, err := h.getMimeType(cmd.Extension)
	if err != nil {
		return "", err
	}
	thumbnailBlurhash, err := h.thumbnailServce.GenerateBlurHash(cmd.File, cmd.Extension)
	if err != nil {
		return "", err
	}

	fileId, err := h.fileRepo.Upload(ctx, cmd.File, mimeType, cmd.Extension)
	if err != nil {
		return "", err
	}

	thumbnailId, err := h.thumbnailRepo.Create(ctx, port.NewThumbnail{
		FileId:   fileId,
		Blurhash: thumbnailBlurhash,
	})
	if err != nil {
		return "", err
	}

	return thumbnailId, nil
}

func (h addThumbnailHandler) getMimeType(extension string) (string, error) {
	switch extension {
	case ".png":
		return "image/png", nil
	case ".jpg", ".jpeg":
		return "image/jpeg", nil
	default:
		return "", errors.New("unsupported image format")
	}
}
