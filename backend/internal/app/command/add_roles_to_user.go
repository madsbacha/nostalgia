package command

import (
	"context"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/core/port"
)

type AddRolesToUser struct {
	UserId string
	Roles  []string
}

type AddRolesToUserHandler decorator.CommandHandler[AddRolesToUser]

type addRolesToUserHandler struct {
	userRepo port.UserRepository
}

func NewAddRolesToUserHandler(userRepo port.UserRepository, logger *logrus.Entry) AddRolesToUserHandler {
	return decorator.ApplyCommandDecorators[AddRolesToUser](
		addRolesToUserHandler{userRepo: userRepo},
		logger,
	)
}

func (h addRolesToUserHandler) Handle(ctx context.Context, cmd AddRolesToUser) error {
	for _, role := range cmd.Roles {
		if err := h.userRepo.AddRoleForUser(ctx, cmd.UserId, role); err != nil {
			return err
		}
	}
	return nil
}
