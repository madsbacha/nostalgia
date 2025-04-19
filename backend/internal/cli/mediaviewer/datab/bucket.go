package datab

import (
	"database/sql"
	"strings"
)

type BucketObject struct {
	Id        int64
	Name      string
	CreatedAt int64
}

type BucketObjectRepo struct {
	Context   *Context
	sqlDb     *sql.DB
	TableName string
}

func (repo *BucketObjectRepo) GetTableName() string {
	return repo.TableName
}

func (repo *BucketObjectRepo) GetSelectProperties() []string {
	return strings.Split("id,name,created_at", ",")
}

func (repo *BucketObjectRepo) Db() *sql.DB {
	return repo.sqlDb
}

func (repo *BucketObjectRepo) Scan(s func(dest ...interface{}) error) (*BucketObject, error) {
	obj := BucketObject{}
	err := s(
		&obj.Id,
		&obj.Name,
		&obj.CreatedAt)

	return &obj, err
}

func GetBucketObjectRepository(context *Context) BucketObjectRepo {
	return BucketObjectRepo{
		Context:   context,
		sqlDb:     context.Db,
		TableName: "files",
	}
}

func (repo *BucketObjectRepo) Insert(obj BucketObject) (sql.Result, error) {
	return Query[*BucketObject](repo).Insert(map[string]interface{}{
		"name":       obj.Name,
		"created_at": obj.CreatedAt,
	})
}

func (repo *BucketObjectRepo) GetById(id int64) (BucketObject, error) {
	obj, err := Query[*BucketObject](repo).Where("id = ?", id).GetFirst()
	if err != nil {
		return BucketObject{}, err
	}
	return *obj, nil
}

func (repo *BucketObjectRepo) GetByKey(key string) (BucketObject, error) {
	obj, err := Query[*BucketObject](repo).Where("name = ?", key).GetFirst()
	if err != nil {
		return BucketObject{}, err
	}
	return *obj, nil
}
