package web

import (
	"context"
	"io"
	"net/http"
	"nostalgia/internal/app"
	"nostalgia/internal/app/command"
	"nostalgia/internal/app/request"
	"nostalgia/internal/core/domain"
)

func Initialize(ctx context.Context, app app.Application) error {
	imageReader, extension, err := getDefaultThumbnailImage(ctx, app)
	if err != nil {
		return err
	}

	thumbnailId, err := app.Requests.AddThumbnail.Handle(ctx, request.AddThumbnail{
		File:      imageReader,
		Extension: extension,
	})
	if err != nil {
		return err
	}

	err = app.Commands.SetDefaultThumbnail.Handle(ctx, command.SetDefaultThumbnail{
		ThumbnailId: thumbnailId,
	})
	if err != nil {
		return err
	}

	err = app.Commands.SetSetting.Handle(ctx, command.SetSetting{
		Key:   domain.SettingTitle,
		Value: "nostalgia",
	})
	if err != nil {
		return err
	}

	return app.Commands.SetSetting.Handle(ctx, command.SetSetting{
		Key:   domain.SettingIsInitialized,
		Value: "true",
	})
}

func getDefaultThumbnailImage(ctx context.Context, app app.Application) (io.Reader, string, error) {
	res, err := http.Get("https://images.unsplash.com/photo-1588345921523-c2dcdb7f1dcd?w=800&dpr=2&q=80")
	if err != nil {
		return nil, "", err
	}
	if res.StatusCode != http.StatusOK {
		return nil, "", err
	}
	return res.Body, ".jpg", nil
}
