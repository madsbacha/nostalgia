package request

import (
	"context"
	"github.com/sirupsen/logrus"
	"io"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/common/errors"
	"nostalgia/internal/core/port"
	"slices"
	"strings"
	"time"
)

type AddMedia struct {
	UserId      string
	File        io.Reader
	Extension   string
	MimeType    string
	Title       string
	Description string
	Tags        []string
}

type AddMediaHandler decorator.RequestHandler[AddMedia, string]

type addMediaHandler struct {
	mediaRepo        port.MediaRepository
	userRepo         port.UserRepository
	fileRepo         port.FileRepository
	settingRepo      port.SettingRepository
	thumbnailRepo    port.ThumbnailRepository
	thumbnailService port.ThumbnailService
}

func NewAddMediaHandler(mediaRepo port.MediaRepository, userRepo port.UserRepository, fileRepo port.FileRepository, settingRepo port.SettingRepository, thumbnailRepo port.ThumbnailRepository, thumbnailService port.ThumbnailService, logger *logrus.Entry) AddMediaHandler {
	return decorator.ApplyRequestDecorators[AddMedia, string](
		addMediaHandler{
			mediaRepo:        mediaRepo,
			userRepo:         userRepo,
			fileRepo:         fileRepo,
			settingRepo:      settingRepo,
			thumbnailRepo:    thumbnailRepo,
			thumbnailService: thumbnailService},
		logger,
	)
}

func (h addMediaHandler) Handle(ctx context.Context, cmd AddMedia) (string, error) {

	// 1. Check if the user exists
	exists, err := h.userRepo.ExistsById(ctx, cmd.UserId)
	if err != nil {
		return "", errors.NewSlugError(err.Error(), "unable-to-check-user-exists")
	}
	if !exists {
		return "", errors.NewSlugError("user does not exist", "user-not-found")
	}

	// 2. Check if the file, filetype and extension is valid
	file := cmd.File
	if file == nil {
		return "", errors.NewSlugError("file is nil", "file-not-found")
	}
	if cmd.Extension == "" {
		return "", errors.NewSlugError("file extension is empty", "file-extension-empty")
	} else if cmd.Extension != ".mp4" && cmd.Extension != ".mov" && cmd.Extension != ".webm" {
		return "", errors.NewSlugError("file extension is not supported", "file-extension-not-supported")
	}

	// 3. Check if the title is valid
	if cmd.Title == "" {
		return "", errors.NewSlugError("title is empty", "title-empty")
	}

	if len(cmd.Description) > 1000 {
		return "", errors.NewSlugError("description is too long", "description-too-long")
	}

	if !h.validExtension(cmd.Extension) {
		return "", errors.NewSlugError("invalid file extension", "invalid-file-extension")
	}

	// 8. Save the file to the server
	fileId, err := h.fileRepo.Upload(ctx, file, cmd.MimeType, cmd.Extension)
	if err != nil {
		return "", errors.NewSlugError(err.Error(), "unable-to-upload-file")
	}

	fileOnFile, err := h.fileRepo.GetById(ctx, fileId)
	if err != nil {
		return "", errors.NewSlugError(err.Error(), "unable-to-get-file")
	}

	generatedThumbnail, err := h.thumbnailService.GenerateFromVideo(ctx, fileOnFile.Path)
	if err != nil {
		return "", errors.NewSlugError(err.Error(), "unable-to-generate-thumbnail")
	}

	thumbnailFileId, err := h.fileRepo.Upload(ctx, generatedThumbnail.File, generatedThumbnail.MimeType, generatedThumbnail.Extension)
	if err != nil {
		return "", errors.NewSlugError(err.Error(), "unable-to-upload-thumbnail")
	}

	thumbnailId, err := h.thumbnailRepo.Create(ctx, port.NewThumbnail{
		FileId:   thumbnailFileId,
		Blurhash: generatedThumbnail.Blurhash,
	})
	if err != nil {
		return "", err
	}

	uploadedAt := time.Now()

	// 9. Save the media to the database
	mediaId, err := h.mediaRepo.AddMedia(ctx, fileId, thumbnailId, cmd.Title, cmd.Description, cmd.UserId, cmd.Tags, uploadedAt)
	if err != nil {
		return "", errors.NewSlugError(err.Error(), "unable-to-add-media")
	}
	return mediaId, nil
}

func (h addMediaHandler) validExtension(ext string) bool {
	ext = strings.ToLower(ext)
	ext = strings.TrimPrefix(ext, ".")
	validExtensions := []string{
		"mp4",
		"mov",
		"webm",
	}
	// TODO: Align extensions across frontend/backend
	return slices.Contains(validExtensions, ext)
}
