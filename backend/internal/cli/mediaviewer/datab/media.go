package datab

import (
	"database/sql"
	"strings"
)

type Media struct {
	Id           string         `json:"id"`
	Title        string         `json:"title"`
	ObjectId     int64          `json:"object_id"`
	UploadedAt   int64          `json:"uploaded_at"`
	Visibility   int            `json:"visibility"`
	Views        int64          `json:"views"`
	UploadedBy   int64          `json:"uploaded_by"`
	Description  string         `json:"description"`
	ThumbnailKey sql.NullString `json:"thumbnail_key"`
	Tags         []string       `json:"tags"`
	Deleted      bool           `json:"deleted"`
}

type MediaRepo struct {
	Context   *Context
	sqlDb     *sql.DB
	TableName string
}

type MediaQueryResult struct {
	items []Media
}

func (repo *MediaRepo) GetTableName() string {
	return repo.TableName
}

func (repo *MediaRepo) GetSelectProperties() []string {
	return strings.Split("id,title,object_id,uploaded_at,visibility,views,uploaded_by,description,thumbnail_key,tags,deleted", ",")
}

func (repo *MediaRepo) Db() *sql.DB {
	return repo.sqlDb
}

func ParseTags(rawTags string) []string {
	tags := strings.Split(strings.ToLower(rawTags), ",")
	if tags == nil {
		return []string{}
	}
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}
	return tags
}

func (repo *MediaRepo) Scan(s func(dest ...interface{}) error) (*Media, error) {
	media := Media{}
	var rawTags string
	err := s(
		&media.Id,
		&media.Title,
		&media.ObjectId,
		&media.UploadedAt,
		&media.Visibility,
		&media.Views,
		&media.UploadedBy,
		&media.Description,
		&media.ThumbnailKey,
		&rawTags,
		&media.Deleted)
	if err != nil {
		return &media, err
	}

	media.Tags = ParseTags(rawTags)

	return &media, nil
}

func GetMediaRepository(context *Context) MediaRepo {
	return MediaRepo{
		Context:   context,
		sqlDb:     context.Db,
		TableName: "media",
	}
}

func (repo *MediaRepo) Insert(media Media) (string, sql.Result, error) {
	panic("implement me")
}

func (repo *MediaRepo) GetById(id string) (*Media, error) {
	return Query[*Media](repo).Where("id = ? AND deleted = 0", id).GetFirst()
}

func (repo *MediaRepo) GetAll() ([]*Media, error) {
	media, err := Query[*Media](repo).Where("deleted = 0").GetAll()
	if err != nil {
		return nil, err
	}

	return media, nil
}

func (repo *MediaRepo) GetAllWithoutThumbnail() ([]*Media, error) {
	return Query[*Media](repo).Where("deleted = 0 AND (thumbnail_key IS NULL OR thumbnail_key == \"\")").GetAll()
}

func (repo *MediaRepo) GetForUser(id int64) ([]*Media, error) {
	return Query[*Media](repo).Where("(uploaded_by = ? OR visibility > 0) AND deleted = 0", id).GetAll()
}

func (repo *MediaRepo) Update(media *Media) (sql.Result, error) {
	return Query[*Media](repo).Where("id = ? AND deleted = 0", media.Id).Update(map[string]interface{}{
		"title":       media.Title,
		"uploaded_at": media.UploadedAt,
		"visibility":  media.Visibility,
		"uploaded_by": media.UploadedBy,
		"description": media.Description,
		"tags":        strings.Join(media.Tags, ","),
	})
}

func (repo *MediaRepo) DeleteById(id string) (sql.Result, error) {
	return Query[*Media](repo).Where("id = ?", id).Update(map[string]interface{}{
		"deleted": 1,
	})
}

func (repo *MediaRepo) PlusOneView(id string) error {
	_, err := Query[*Media](repo).Where("id = ?", id).AddToColumn("views", 1)
	return err
}

func (repo *MediaRepo) GetAllTags() ([]string, error) {
	rawTags, err := Query[*Media](repo).Where("deleted = 0").Column("tags").GetAllStrings()
	if err != nil {
		return nil, err
	}
	tags := make([]string, 0)
	for _, rawTag := range rawTags {
		if strings.TrimSpace(rawTag) == "" {
			continue
		}
		rowTags := ParseTags(rawTag)
		for _, tag := range rowTags {
			newTag := strings.TrimSpace(tag)
			exists := false
			for _, existingTag := range tags {
				if strings.ToLower(newTag) == strings.ToLower(existingTag) {
					exists = true
					break
				}
			}
			if !exists {
				tags = append(tags, newTag)
			}
		}
	}

	return tags, nil
}

func (repo *MediaRepo) SetThumbnail(id, key string) error {
	_, err := Query[*Media](repo).Where("id = ?", id).Update(map[string]interface{}{
		"thumbnail_key": key,
	})

	return err
}

func (repo *MediaRepo) GetByObjId(id int64) (*Media, error) {
	return Query[*Media](repo).Where("object_id = ?", id).GetFirst()
}
