package query

import (
	"context"
	"github.com/sirupsen/logrus"
	"nostalgia/internal/common/decorator"
	"nostalgia/internal/core/port"
)

type GetRolesForUser struct {
	UserId string
}

type GetRolesForUserHandler decorator.QueryHandler[GetRolesForUser, []string]

type getRolesForUserHandler struct {
	userRepo port.UserRepository
}

func NewGetRolesForUserHandler(userRepo port.UserRepository, logger *logrus.Entry) GetRolesForUserHandler {
	return decorator.ApplyQueryDecorators[GetRolesForUser, []string](
		getRolesForUserHandler{userRepo: userRepo},
		logger,
	)
}

func (h getRolesForUserHandler) Handle(ctx context.Context, query GetRolesForUser) ([]string, error) {
	return h.userRepo.GetRolesForUser(ctx, query.UserId)
}
