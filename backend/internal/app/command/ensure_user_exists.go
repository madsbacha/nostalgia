package command

import (
	"context"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/common/errors"
	"nostalgia/internal/core/port"
)

type EnsureUserExists struct {
	DiscordId string
	Username  string
}

type EnsureUserExistsHandler decorator.CommandHandler[EnsureUserExists]

type ensureUserExistsHandler struct {
	userRepo port.UserRepository
}

func NewEnsureUserExistsHandler(userRepo port.UserRepository, logger *logrus.Entry) EnsureUserExistsHandler {
	return decorator.ApplyCommandDecorators[EnsureUserExists](
		ensureUserExistsHandler{userRepo: userRepo},
		logger,
	)
}

func (h ensureUserExistsHandler) Handle(ctx context.Context, cmd EnsureUserExists) error {
	exists, err := h.userRepo.ExistsByDiscordId(ctx, cmd.DiscordId)
	if err != nil {
		return err
	}

	if !exists {
		err := h.userRepo.Insert(ctx, port.NewUser{
			DiscordId: cmd.DiscordId,
			Username:  cmd.Username,
		})
		if err != nil {
			return errors.NewSlugError(err.Error(), "could-not-insert-user")
		}
	}

	return nil
}
