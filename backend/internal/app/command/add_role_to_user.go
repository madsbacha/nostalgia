package command

import (
	"context"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/core/port"
)

type AddRoleToUser struct {
	UserId string
	Role   string
}

type AddRoleToUserHandler decorator.CommandHandler[AddRoleToUser]

type addRoleToUserHandler struct {
	userRepo port.UserRepository
}

func NewAddRoleToUserHandler(userRepo port.UserRepository, logger *logrus.Entry) AddRoleToUserHandler {
	return decorator.ApplyCommandDecorators[AddRoleToUser](
		addRoleToUserHandler{userRepo: userRepo},
		logger,
	)
}

func (h addRoleToUserHandler) Handle(ctx context.Context, cmd AddRoleToUser) error {
	return h.userRepo.AddRoleForUser(ctx, cmd.UserId, cmd.Role)
}
