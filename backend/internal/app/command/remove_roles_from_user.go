package command

import (
	"context"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/core/port"
)

type RemoveRolesFromUser struct {
	UserId string
	Roles  []string
}

type RemoveRolesFromUserHandler decorator.CommandHandler[RemoveRolesFromUser]

type removeRolesFromUserHandler struct {
	userRepo port.UserRepository
}

func NewRemoveRolesFromUserHandler(userRepo port.UserRepository, logger *logrus.Entry) RemoveRolesFromUserHandler {
	return decorator.ApplyCommandDecorators[RemoveRolesFromUser](
		removeRolesFromUserHandler{userRepo: userRepo},
		logger,
	)
}

func (h removeRolesFromUserHandler) Handle(ctx context.Context, cmd RemoveRolesFromUser) error {
	for _, role := range cmd.Roles {
		if err := h.userRepo.RemoveRoleForUser(ctx, cmd.UserId, role); err != nil {
			return err
		}
	}
	return nil
}
