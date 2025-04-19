package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"io"
	"nostalgia/internal/core/port"
	"nostalgia/pkg/dataquery"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type File struct {
	Id       int64
	Path     string
	MimeType string
}

const (
	TableFile          = "file"
	FileColumnId       = "id"
	FileColumnPath     = "path"
	FileColumnMimeType = "mime_type"
)

type FileScanner struct{}

func NewFileScanner() *FileScanner {
	return &FileScanner{}
}

func (s FileScanner) GetSelectProperties() []string {
	return []string{
		FileColumnId,
		FileColumnPath,
		FileColumnMimeType,
	}
}

func (s FileScanner) Scan(sFunc func(dest ...interface{}) error) (*File, error) {
	file := File{}
	err := sFunc(
		&file.Id,
		&file.Path,
		&file.MimeType)
	return &file, err
}

func (r FileLocalRepository) query(ctx context.Context) *dataquery.QueryBuilder[*File] {
	scanner := NewFileScanner()
	return dataquery.QueryContext[*File](ctx, r.db, TableFile, scanner)
}

type FileLocalRepository struct {
	db       *sql.DB
	basePath string
}

func NewFileLocalRepository(db *sql.DB, basePath string) *FileLocalRepository {
	return &FileLocalRepository{
		db:       db,
		basePath: basePath,
	}
}

func (r FileLocalRepository) Upload(ctx context.Context, file io.Reader, mimeType string, extension string) (string, error) {
	ext := strings.ToLower(extension)
	ext = strings.TrimPrefix(ext, ".")

	uid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	filename := uid.String() + "." + ext
	fullFilepath := filepath.Join(r.basePath, filename)
	f, err := os.Create(fullFilepath)
	if err != nil {
		return "", err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	_, err = io.Copy(f, file)
	if err != nil {
		return "", err
	}

	// TODO: Escape filename and mimetype
	result, err := r.query(ctx).Insert(map[string]interface{}{
		FileColumnPath:     filename,
		FileColumnMimeType: mimeType,
	})
	if err != nil {
		return "", err
	}

	record, err := result.LastInsertId()

	return strconv.FormatInt(record, 10), err
}

func (r FileLocalRepository) GetById(ctx context.Context, id string) (*port.File, error) {
	where := fmt.Sprintf("%s = ?", FileColumnId)
	file, err := r.query(ctx).Where(where, id).First()
	if err != nil {
		return nil, err
	}

	return &port.File{
		Id:       strconv.FormatInt(file.Id, 10),
		Path:     filepath.Join(r.basePath, file.Path),
		MimeType: file.MimeType,
	}, nil
}
