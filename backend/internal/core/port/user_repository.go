package port

import (
	"context"
	"nostalgia/internal/core/domain"
)

type NewUser struct {
	DiscordId string
	Username  string
	Avatar    string
	Roles     []string
}

type UserRepository interface {
	Insert(ctx context.Context, user NewUser) error
	ExistsByDiscordId(ctx context.Context, discordId string) (bool, error)
	ExistsById(ctx context.Context, id string) (bool, error)
	GetByDiscordId(ctx context.Context, discordId string) (*domain.User, error)
	GetById(ctx context.Context, id string) (*domain.User, error)
	Update(ctx context.Context, user domain.User) error
	GetRolesForUser(ctx context.Context, userId string) ([]string, error)
	AddRoleForUser(ctx context.Context, userId string, role string) error
	RemoveRoleForUser(ctx context.Context, userId string, role string) error
	GetAll(ctx context.Context) ([]domain.User, error)
}
