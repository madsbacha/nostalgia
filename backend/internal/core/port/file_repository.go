package port

import (
	"context"
	"io"
)

type File struct {
	Id       string
	Path     string
	MimeType string
}

type FileRepository interface {
	Upload(ctx context.Context, file io.Reader, mimeType string, extension string) (string, error)
	GetById(ctx context.Context, id string) (*File, error)
}
