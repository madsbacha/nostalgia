package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"nostalgia/internal/common/util"
	"nostalgia/internal/core/domain"
	"nostalgia/internal/core/port"
	"nostalgia/pkg/dataquery"
	"strconv"
	"strings"
)

type User struct {
	Id        int64
	DiscordId string
	Username  string
	Avatar    sql.NullString
	Roles     string
}

const (
	TableUser           = "user"
	UserColumnId        = "id"
	UserColumnDiscordId = "discord_id"
	UserColumnUsername  = "username"
	UserColumnAvatar    = "avatar"
	UserColumnRoles     = "roles"
)

type UserScanner struct{}

func NewUserScanner() *UserScanner {
	return &UserScanner{}
}

func (scanner UserScanner) GetSelectProperties() []string {
	return []string{
		UserColumnId,
		UserColumnDiscordId,
		UserColumnUsername,
		UserColumnAvatar,
		UserColumnRoles,
	}
}

func (scanner UserScanner) Scan(s func(dest ...interface{}) error) (*User, error) {
	user := User{}
	err := s(
		&user.Id,
		&user.DiscordId,
		&user.Username,
		&user.Avatar,
		&user.Roles)
	return &user, err
}

func (r UserSqliteRepository) query(ctx context.Context) *dataquery.QueryBuilder[*User] {
	scanner := NewUserScanner()
	return dataquery.QueryContext[*User](ctx, r.db, TableUser, scanner)
}

func NewUserSqliteRepository(db *sql.DB) *UserSqliteRepository {
	return &UserSqliteRepository{
		db: db,
	}
}

type UserSqliteRepository struct {
	db *sql.DB
}

func (r UserSqliteRepository) Update(ctx context.Context, user domain.User) error {
	where := fmt.Sprintf("%s = ?", UserColumnId)
	_, err := r.query(ctx).Where(where, user.Id).Update(map[string]interface{}{
		UserColumnUsername: user.Username,
		UserColumnAvatar:   user.Avatar,
	})
	return err
}

func (r UserSqliteRepository) Insert(ctx context.Context, user port.NewUser) error {
	_, err := r.query(ctx).Insert(map[string]interface{}{
		UserColumnDiscordId: user.DiscordId,
		UserColumnUsername:  user.Username,
		UserColumnAvatar:    user.Avatar,
		UserColumnRoles:     strings.Join(user.Roles, ","),
	})
	return err
}

func (r UserSqliteRepository) ExistsByDiscordId(ctx context.Context, discordId string) (bool, error) {
	where := fmt.Sprintf("%s = ?", UserColumnDiscordId)
	exists, err := r.query(ctx).Where(where, discordId).Exists()
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r UserSqliteRepository) ExistsById(ctx context.Context, id string) (bool, error) {
	where := fmt.Sprintf("%s = ?", UserColumnId)
	exists, err := r.query(ctx).Where(where, id).Exists()
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r UserSqliteRepository) GetByDiscordId(ctx context.Context, discordId string) (*domain.User, error) {
	where := fmt.Sprintf("%s = ?", UserColumnDiscordId)
	user, err := r.query(ctx).Where(where, discordId).First()
	if err != nil {
		return nil, err
	}
	return &domain.User{
		Id:        strconv.FormatInt(user.Id, 10),
		DiscordId: user.DiscordId,
		Username:  user.Username,
		Avatar:    user.Avatar.String,
		Roles:     strings.Split(user.Roles, ","),
	}, nil
}

func (r UserSqliteRepository) GetById(ctx context.Context, id string) (*domain.User, error) {
	where := fmt.Sprintf("%s = ?", UserColumnId)
	user, err := r.query(ctx).Where(where, id).First()
	if err != nil {
		return nil, err
	}
	return &domain.User{
		Id:        strconv.FormatInt(user.Id, 10),
		DiscordId: user.DiscordId,
		Username:  user.Username,
		Avatar:    user.Avatar.String,
		Roles:     strings.Split(user.Roles, ","),
	}, nil
}

func (r UserSqliteRepository) GetRolesForUser(ctx context.Context, userId string) ([]string, error) {
	where := fmt.Sprintf("%s = ?", UserColumnId)
	user, err := r.query(ctx).Where(where, userId).First()
	if err != nil {
		return nil, err
	}

	roles := strings.Split(user.Roles, ",")

	return roles, nil
}

func (r UserSqliteRepository) AddRoleForUser(ctx context.Context, userId string, role string) error {
	where := fmt.Sprintf("%s = ?", UserColumnId)
	user, err := r.query(ctx).Where(where, userId).First()
	if err != nil {
		return err
	}

	roles := strings.Split(user.Roles, ",")
	if !util.Contains[string](roles, role) {
		roles = append(roles, role)
	}
	roles = util.RemoveEmptyTags(roles)

	_, err = r.query(ctx).Where(where, userId).Update(map[string]interface{}{
		UserColumnRoles: strings.Join(roles, ","),
	})
	return err
}

func (r UserSqliteRepository) RemoveRoleForUser(ctx context.Context, userId string, role string) error {
	where := fmt.Sprintf("%s = ?", UserColumnId)
	user, err := r.query(ctx).Where(where, userId).First()
	if err != nil {
		return err
	}

	roles := strings.Split(user.Roles, ",")
	roles = util.Remove(roles, role)

	_, err = r.query(ctx).Where(where, userId).Update(map[string]interface{}{
		UserColumnRoles: strings.Join(roles, ","),
	})
	return err
}

func (r UserSqliteRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	users, err := r.query(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	domainUsers := make([]domain.User, len(users))
	for i, user := range users {
		domainUsers[i] = domain.User{
			Id:        strconv.FormatInt(user.Id, 10),
			DiscordId: user.DiscordId,
			Username:  user.Username,
			Avatar:    user.Avatar.String,
			Roles:     strings.Split(user.Roles, ","),
		}
	}

	return domainUsers, nil
}
