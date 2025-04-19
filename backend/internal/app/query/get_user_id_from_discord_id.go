package query

import (
	"context"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/common/errors"
	"nostalgia/internal/core/port"
)

type GetUserIdFromDiscordId struct {
	DiscordId string
}

type GetUserIdFromDiscordIdHandler decorator.QueryHandler[GetUserIdFromDiscordId, string]

type getUserIdFromDiscordIdHandler struct {
	userRepo port.UserRepository
}

func NewGetUserIdFromDiscordId(userRepo port.UserRepository, logger *logrus.Entry) GetUserIdFromDiscordIdHandler {
	return decorator.ApplyQueryDecorators[GetUserIdFromDiscordId, string](
		getUserIdFromDiscordIdHandler{userRepo: userRepo},
		logger,
	)
}

func (h getUserIdFromDiscordIdHandler) Handle(ctx context.Context, q GetUserIdFromDiscordId) (string, error) {
	user, err := h.userRepo.GetByDiscordId(ctx, q.DiscordId)
	if err != nil {
		return "", errors.NewSlugError(err.Error(), "could-not-get-user-by-discord-id")
	}

	return user.Id, nil
}
