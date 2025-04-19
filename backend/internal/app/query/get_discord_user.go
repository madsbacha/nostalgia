package query

import (
	"context"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/common/errors"
	"nostalgia/pkg/discord"
	"nostalgia/pkg/oauth2"
)

type GetDiscordUser struct {
	Token oauth2.Token
}

type GetDiscordUserHandler decorator.QueryHandler[GetDiscordUser, *discord.User]

type getDiscordUserHandler struct {
	discord discord.Client
}

func NewGetDiscordUserHandler(discordCtx discord.Client, logger *logrus.Entry) GetDiscordUserHandler {
	return decorator.ApplyQueryDecorators[GetDiscordUser, *discord.User](
		getDiscordUserHandler{discord: discordCtx},
		logger,
	)
}

func (s getDiscordUserHandler) Handle(ctx context.Context, q GetDiscordUser) (*discord.User, error) {
	user, err := s.discord.GetUser(q.Token)
	if err != nil {
		return nil, errors.NewSlugError(err.Error(), "could-not-get-discord-user")
	}

	return &user, nil
}
