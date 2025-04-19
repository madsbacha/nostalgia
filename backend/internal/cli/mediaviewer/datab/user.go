package datab

import (
	"database/sql"
	"encoding/json"
	"strings"
)

type User struct {
	Id        int64
	DiscordId string         `json:"id"`
	Username  string         `json:"username"`
	Avatar    string         `json:"avatar"`
	Admin     bool           `json:"admin"`
	UploadKey sql.NullString `json:"upload_key"`
}

type UserRepo struct {
	Context   *Context
	sqlDb     *sql.DB
	TableName string
}

type UserQueryResult struct {
	items []User
}

func (repo *UserRepo) GetTableName() string {
	return repo.TableName
}

func (repo *UserRepo) GetSelectProperties() []string {
	return strings.Split("id,discord_id,username,avatar,upload_key,admin", ",")
}

func (repo *UserRepo) Db() *sql.DB {
	return repo.sqlDb
}

func (repo *UserRepo) Scan(s func(dest ...interface{}) error) (*User, error) {
	user := User{}
	err := s(
		&user.Id,
		&user.DiscordId,
		&user.Username,
		&user.Avatar,
		&user.UploadKey,
		&user.Admin)
	return &user, err
}

func GetUserRepository(context *Context) UserRepo {
	return UserRepo{
		Context:   context,
		sqlDb:     context.Db,
		TableName: "users",
	}
}

func UserFromJson(s []byte) (User, error) {
	user := User{}

	err := json.Unmarshal(s, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (repo *UserRepo) Insert(user *User) (sql.Result, error) {
	panic("not implemented")
}

func (repo *UserRepo) Exists(discordId string) (bool, error) {
	return Query[*User](repo).Where("discord_id = ?", discordId).Exists()
}

func (repo *UserRepo) Update(user *User) (sql.Result, error) {
	panic("not implemented")
}

func (repo *UserRepo) InsertOrUpdate(user *User) (sql.Result, error) {
	exists, err := repo.Exists(user.DiscordId)
	if err != nil {
		return nil, err
	}

	if exists {
		if !user.UploadKey.Valid {
			existingUser, err := repo.GetByDiscordId(user.DiscordId)
			if err != nil {
				return nil, err
			}

			if existingUser.UploadKey.Valid {
				user.UploadKey = existingUser.UploadKey
			}
		}

		return repo.Update(user)
	} else {
		return repo.Insert(user)
	}
}

func (repo *UserRepo) GetByDiscordId(id string) (*User, error) {
	return Query[*User](repo).Where("discord_id = ?", id).GetFirst()
}

func (repo *UserRepo) GetById(id int64) (*User, error) {
	return Query[*User](repo).Where("id = ?", id).GetFirst()
}

func (repo *UserRepo) GetByUploadKey(id string) (*User, error) {
	return Query[*User](repo).Where("upload_key = ?", id).GetFirst()
}
