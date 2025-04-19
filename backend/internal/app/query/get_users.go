package query

import (
	"context"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/core/domain"
	"nostalgia/internal/core/port"
)

type GetUsers struct{}

type GetUsersHandler decorator.QueryHandler[GetUsers, []domain.User]

type getUsersHandler struct {
	userRepo port.UserRepository
}

func NewGetUsersHandler(userRepo port.UserRepository, logger *logrus.Entry) GetUsersHandler {
	return decorator.ApplyQueryDecorators[GetUsers, []domain.User](
		getUsersHandler{
			userRepo: userRepo,
		},
		logger,
	)
}

func (h getUsersHandler) Handle(ctx context.Context, query GetUsers) ([]domain.User, error) {
	users, err := h.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}
