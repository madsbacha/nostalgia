package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"nostalgia/internal/core/port"
	"nostalgia/pkg/dataquery"
)

type Thumbnail struct {
	Id       string
	FileId   string
	BlurHash string
}

const (
	TableThumbnail          = "thumbnail"
	ThumbnailColumnId       = "id"
	ThumbnailColumnFileId   = "file_id"
	ThumbnailColumnBlurhash = "blurhash"
)

type ThumbnailScanner struct{}

func NewThumbnailScanner() *ThumbnailScanner {
	return &ThumbnailScanner{}
}

func (s ThumbnailScanner) GetSelectProperties() []string {
	return []string{
		ThumbnailColumnId,
		ThumbnailColumnFileId,
		ThumbnailColumnBlurhash,
	}
}

func (s ThumbnailScanner) Scan(sFunc func(dest ...interface{}) error) (*Thumbnail, error) {
	thumbnail := Thumbnail{}
	err := sFunc(
		&thumbnail.Id,
		&thumbnail.FileId,
		&thumbnail.BlurHash)
	return &thumbnail, err
}

func (r ThumbnailLocalRepository) query(ctx context.Context) *dataquery.QueryBuilder[*Thumbnail] {
	scanner := NewThumbnailScanner()
	return dataquery.QueryContext[*Thumbnail](ctx, r.db, TableThumbnail, scanner)
}

type ThumbnailLocalRepository struct {
	db *sql.DB
}

func NewThumbnailSqliteRepository(db *sql.DB) *ThumbnailLocalRepository {
	return &ThumbnailLocalRepository{
		db: db,
	}
}

func (r ThumbnailLocalRepository) GetById(ctx context.Context, id string) (*port.Thumbnail, error) {
	where := fmt.Sprintf("%s = ?", ThumbnailColumnId)
	thumbnail, err := r.query(ctx).Where(where, id).First()
	if err != nil {
		return nil, err
	}

	return &port.Thumbnail{
		Id:       thumbnail.Id,
		FileId:   thumbnail.FileId,
		Blurhash: thumbnail.BlurHash,
	}, nil
}

func (r ThumbnailLocalRepository) Create(ctx context.Context, thumbnail port.NewThumbnail) (string, error) {
	result, err := r.query(ctx).Insert(map[string]interface{}{
		ThumbnailColumnFileId:   thumbnail.FileId,
		ThumbnailColumnBlurhash: thumbnail.Blurhash,
	})
	if err != nil {
		return "", err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", id), nil
}
