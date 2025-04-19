package command

import (
	"context"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/core/port"
)

type UpdateUser struct {
	Id     string
	Avatar string
}

type UpdateUserHandler decorator.CommandHandler[UpdateUser]

type updateUserHandler struct {
	userRepo port.UserRepository
}

func NewUpdateUserHandler(userRepo port.UserRepository, logger *logrus.Entry) UpdateUserHandler {
	return decorator.ApplyCommandDecorators[UpdateUser](
		updateUserHandler{userRepo: userRepo},
		logger,
	)
}

func (h updateUserHandler) Handle(ctx context.Context, cmd UpdateUser) error {
	user, err := h.userRepo.GetById(ctx, cmd.Id)
	if err != nil {
		return err
	}
	user.Avatar = cmd.Avatar

	return h.userRepo.Update(ctx, *user)
}
