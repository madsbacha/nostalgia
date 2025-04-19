package query

import (
	"context"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/core/domain"
	"nostalgia/internal/core/port"
)

type GetUserById struct {
	Id string
}

type GetUserByIdHandler decorator.QueryHandler[GetUserById, *domain.User]

type getUserByIdHandler struct {
	userRepo port.UserRepository
}

func NewGetUserById(userRepo port.UserRepository, logger *logrus.Entry) GetUserByIdHandler {
	return decorator.ApplyQueryDecorators[GetUserById, *domain.User](
		getUserByIdHandler{userRepo: userRepo},
		logger,
	)
}

func (h getUserByIdHandler) Handle(ctx context.Context, q GetUserById) (*domain.User, error) {
	user, err := h.userRepo.GetById(ctx, q.Id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
