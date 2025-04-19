package query

import (
	"context"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/common/errors"
	"nostalgia/pkg/discord"
	"nostalgia/pkg/oauth2"
)

type ExchangeDiscordCode struct {
	Code        string
	RedirectUri string
}

type ExchangeDiscordCodeHandler decorator.QueryHandler[ExchangeDiscordCode, *oauth2.Token]

type exchangeDiscordCode struct {
	discordCtx discord.Client
}

func NewExchangeDiscordCode(discordCtx discord.Client, logger *logrus.Entry) ExchangeDiscordCodeHandler {
	return decorator.ApplyQueryDecorators[ExchangeDiscordCode, *oauth2.Token](
		exchangeDiscordCode{discordCtx: discordCtx},
		logger,
	)
}

func (h exchangeDiscordCode) Handle(ctx context.Context, cmd ExchangeDiscordCode) (*oauth2.Token, error) {
	discordToken, err := oauth2.ExchangeCode(
		cmd.Code,
		h.discordCtx.ClientId,
		h.discordCtx.ClientSecret,
		cmd.RedirectUri,
		"https://discord.com/api/oauth2/token")
	if err != nil {
		return nil, errors.NewSlugError(err.Error(), "unable-to-verify-code")
	}

	return discordToken, nil
}
