package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"nostalgia/internal/common/util"
	"nostalgia/internal/core/domain"
	"nostalgia/pkg/dataquery"
	"nostalgia/pkg/utils"
	"strings"
	"time"
)

type Media struct {
	Id          string
	Title       string
	FileId      string
	ThumbnailId string
	UploadedBy  string
	Description string
	Tags        string
	UploadedAt  int64
}

const (
	TableMedia             = "media"
	MediaColumnId          = "id"
	MediaColumnTitle       = "title"
	MediaColumnFileId      = "file_id"
	MediaColumnThumbnailId = "thumbnail_id"
	MediaColumnUploadedBy  = "uploaded_by"
	MediaColumnDescription = "description"
	MediaColumnTags        = "tags"
	MediaColumnUploadedAt  = "uploaded_at"
)

type MediaScanner struct{}

func NewMediaScanner() *MediaScanner {
	return &MediaScanner{}
}

func (scanner MediaScanner) GetSelectProperties() []string {
	return []string{
		MediaColumnId,
		MediaColumnTitle,
		MediaColumnFileId,
		MediaColumnThumbnailId,
		MediaColumnUploadedBy,
		MediaColumnDescription,
		MediaColumnTags,
		MediaColumnUploadedAt,
	}
}

func (scanner MediaScanner) Scan(s func(dest ...interface{}) error) (*Media, error) {
	media := Media{}
	err := s(
		&media.Id,
		&media.Title,
		&media.FileId,
		&media.ThumbnailId,
		&media.UploadedBy,
		&media.Description,
		&media.Tags,
		&media.UploadedAt)
	return &media, err
}

func (m Media) ToDomain() domain.Media {
	return domain.Media{
		Id:          m.Id,
		Title:       m.Title,
		Description: m.Description,
		FileId:      m.FileId,
		ThumbnailId: m.ThumbnailId,
		UploadedBy:  m.UploadedBy,
		Tags:        util.ParseTags(m.Tags),
		UploadedAt:  m.UploadedAt,
	}
}

func (r MediaSqliteRepository) query(ctx context.Context) *dataquery.QueryBuilder[*Media] {
	scanner := NewMediaScanner()
	return dataquery.QueryContext[*Media](ctx, r.db, TableMedia, scanner)
}

func NewMediaSqliteRepository(db *sql.DB) *MediaSqliteRepository {
	return &MediaSqliteRepository{
		db: db,
	}
}

type MediaSqliteRepository struct {
	db *sql.DB
}

func (r MediaSqliteRepository) AddTagToMedia(ctx context.Context, mediaId string, tag string) error {
	where := fmt.Sprintf("%s = ?", MediaColumnId)
	media, err := r.query(ctx).Where(where, mediaId).First()
	if err != nil {
		return err
	}

	tag = strings.TrimSpace(tag)
	tags := util.ParseTags(media.Tags)
	tagExists := false
	for _, existingTag := range tags {
		if strings.ToLower(existingTag) == strings.ToLower(tag) {
			tagExists = true
		}
	}
	if !tagExists {
		tags = append(tags, tag)
	}

	_, err = r.query(ctx).Where(where, mediaId).Update(map[string]interface{}{
		MediaColumnTags: util.JoinTags(tags),
	})
	return err
}

func (r MediaSqliteRepository) RemoveTagFromMedia(ctx context.Context, mediaId string, tag string) error {
	where := fmt.Sprintf("%s = ?", MediaColumnId)
	media, err := r.query(ctx).Where(where, mediaId).First()
	if err != nil {
		return err
	}

	tags := util.ParseTags(media.Tags)
	newTags := make([]string, 0)
	for _, existingTag := range tags {
		if strings.ToLower(existingTag) == strings.ToLower(tag) {
			continue
		}
		newTags = append(newTags, tag)
	}

	_, err = r.query(ctx).Where(where, mediaId).Update(map[string]interface{}{
		MediaColumnTags: util.JoinTags(newTags),
	})
	return err
}

func (r MediaSqliteRepository) AddMedia(ctx context.Context, fileId string, thumbnailId string, title string, description string, uploadedBy string, tags []string, uploadedAt time.Time) (string, error) {
	id, err := utils.GenerateRandomString(10)
	if err != nil {
		return "", err
	}

	_, err = r.query(ctx).Insert(map[string]interface{}{
		MediaColumnId:          id,
		MediaColumnTitle:       title,
		MediaColumnDescription: description,
		MediaColumnFileId:      fileId,
		MediaColumnUploadedBy:  uploadedBy,
		MediaColumnThumbnailId: thumbnailId,
		MediaColumnTags:        util.JoinTags(tags),
		MediaColumnUploadedAt:  uploadedAt.UTC().Unix(),
	})
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r MediaSqliteRepository) GetTagsForMedia(ctx context.Context, mediaId string) ([]string, error) {
	where := fmt.Sprintf("%s = ?", MediaColumnId)
	media, err := r.query(ctx).Where(where, mediaId).First()
	if err != nil {
		return nil, err
	}

	tags := util.ParseTags(media.Tags)
	return tags, nil
}

func (r MediaSqliteRepository) GetAllTags(ctx context.Context) ([]string, error) {
	mediaList, err := r.query(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	tags := make([]string, 0)
	for _, m := range mediaList {
		mediaTags := util.ParseTags(m.Tags)
		for _, tag := range mediaTags {
			mediaTagExists := false
			for _, existingTag := range tags {
				if existingTag == tag {
					mediaTagExists = true
					break
				}
			}
			if !mediaTagExists {
				tags = append(tags, tag)
			}
		}
	}

	return tags, nil
}

func (r MediaSqliteRepository) GetById(ctx context.Context, id string) (domain.Media, error) {
	where := fmt.Sprintf("%s = ?", MediaColumnId)

	media, err := r.query(ctx).Where(where, id).First()
	if err != nil {
		return domain.Media{}, err
	}
	return media.ToDomain(), nil
}

func (r MediaSqliteRepository) Get(ctx context.Context) ([]domain.Media, error) {
	mediaList, err := r.query(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	media := make([]domain.Media, len(mediaList))
	for i, m := range mediaList {
		media[i] = m.ToDomain()
	}
	return media, nil
}
