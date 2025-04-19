package command

import (
	"context"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/core/port"
)

type RemoveRoleFromUser struct {
	UserId string
	Role   string
}

type RemoveRoleFromUserHandler decorator.CommandHandler[RemoveRoleFromUser]

type removeRoleFromUserHandler struct {
	userRepo port.UserRepository
}

func NewRemoveRoleFromUserHandler(userRepo port.UserRepository, logger *logrus.Entry) RemoveRoleFromUserHandler {
	return decorator.ApplyCommandDecorators[RemoveRoleFromUser](
		removeRoleFromUserHandler{userRepo: userRepo},
		logger,
	)
}

func (h removeRoleFromUserHandler) Handle(ctx context.Context, cmd RemoveRoleFromUser) error {
	return h.userRepo.RemoveRoleForUser(ctx, cmd.UserId, cmd.Role)
}
