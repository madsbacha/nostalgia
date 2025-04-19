package web

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"nostalgia/internal/app"
	"nostalgia/internal/app/query"
	"nostalgia/internal/common/env"
	"nostalgia/internal/common/logs"
	"nostalgia/internal/core/domain"
	"nostalgia/internal/web/handler"
	"nostalgia/internal/web/service"
)

func Run() {
	logs.Init()

	ctx := context.Background()
	application := service.NewApplication(ctx)

	if !IsInitialized(ctx, application) {
		if err := Initialize(ctx, application); err != nil {
			log.Panicln(err)
		}
	}

	logger := logrus.NewEntry(logrus.StandardLogger())
	tokenAuth := jwtauth.New("HS256", []byte(env.MustGet("JWT_SECRET")), nil)
	RunHttpServer(func(router chi.Router) http.Handler {
		server := handler.NewHttpServer(application, logger, tokenAuth)
		return server.RegisterServerHandlers(router)
	})
}

func IsInitialized(ctx context.Context, app app.Application) bool {
	setting, err := app.Queries.GetSetting.Handle(ctx, query.GetSetting{
		Key: domain.SettingIsInitialized,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false
		}
		log.Panicln(err)
	}
	return setting.Value == "true"
}
