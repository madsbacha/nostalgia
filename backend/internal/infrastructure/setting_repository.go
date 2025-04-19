package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"nostalgia/internal/core/domain"
	"nostalgia/pkg/dataquery"
)

type Setting struct {
	Id    int64
	Key   string
	Value string
}

const (
	TableSetting       = "setting"
	SettingColumnId    = "id"
	SettingColumnKey   = "key"
	SettingColumnValue = "value"
)

type SettingScanner struct{}

func NewSettingScanner() *SettingScanner {
	return &SettingScanner{}
}

func (scanner SettingScanner) GetSelectProperties() []string {
	return []string{
		SettingColumnId,
		SettingColumnKey,
		SettingColumnValue,
	}
}

func (scanner SettingScanner) Scan(s func(dest ...interface{}) error) (*Setting, error) {
	setting := Setting{}
	err := s(
		&setting.Id,
		&setting.Key,
		&setting.Value)
	return &setting, err
}

func (r SettingSqliteRepository) query(ctx context.Context) *dataquery.QueryBuilder[*Setting] {
	scanner := NewSettingScanner()
	return dataquery.QueryContext[*Setting](ctx, r.db, TableSetting, scanner)
}

type SettingSqliteRepository struct {
	db *sql.DB
}

func NewSettingSqliteRepository(db *sql.DB) *SettingSqliteRepository {
	return &SettingSqliteRepository{db: db}
}

func (r SettingSqliteRepository) GetByKey(ctx context.Context, key string) (domain.Setting, error) {
	where := fmt.Sprintf("%s = ?", SettingColumnKey)
	setting, err := r.query(ctx).Where(where, key).First()
	if err != nil {
		return domain.Setting{}, err
	}
	return domain.Setting{
		Key:   setting.Key,
		Value: setting.Value,
	}, nil
}

func (r SettingSqliteRepository) Set(ctx context.Context, key, value string) error {
	_, err := r.query(ctx).Insert(map[string]interface{}{
		SettingColumnKey:   key,
		SettingColumnValue: value,
	})
	if err != nil {
		return err
	}
	return nil
}
