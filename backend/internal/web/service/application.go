package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"log"
	"nostalgia/internal/app"
	"nostalgia/internal/app/command"
	"nostalgia/internal/app/query"
	"nostalgia/internal/app/request"
	"nostalgia/internal/common/config"
	"nostalgia/internal/common/env"
	"nostalgia/internal/infrastructure"
	"nostalgia/pkg/discord"
	"os"
	"strings"
)

func NewApplication(ctx context.Context) app.Application {
	jwtConfig := config.NewJwtConfig(env.MustGet("JWT_SECRET"))

	db := infrastructure.NewSqliteDatabase(env.MustGet("SQLITE_PATH"))
	if env.GetBool("SQLITE_RUN_MIGRATIONS", false) {
		err := infrastructure.MigrateSqliteDatabase(db)
		if err != nil {
			log.Panicln(err)
		}
	}

	logger := logrus.NewEntry(logrus.StandardLogger())

	settingRepository := infrastructure.NewSettingSqliteRepository(db)
	mediaRepository := infrastructure.NewMediaSqliteRepository(db)
	userRepository := infrastructure.NewUserSqliteRepository(db)
	thumbnailRepository := infrastructure.NewThumbnailSqliteRepository(db)
	fileRepository := infrastructure.NewFileLocalRepository(db, env.MustGet("STORAGE_PATH"))
	thumbnailService := infrastructure.NewThumbnailService(env.MustGet("FFMPEG_PATH"), env.MustGet("TEMP_DATA_PATH"), logger)

	discordCtx := discord.NewClient(
		os.Getenv("DISCORD_CLIENT_ID"),
		os.Getenv("DISCORD_CLIENT_SECRET"),
		strings.Split(os.Getenv("DISCORD_CLIENT_SCOPES"), ","))

	return app.Application{
		Commands: app.Commands{
			AddTagToMedia:       command.NewAddTagToMediaHandler(mediaRepository, logger),
			RemoveTagFromMedia:  command.NewRemoveTagFromMediaHandler(mediaRepository, logger),
			EnsureUserExists:    command.NewEnsureUserExistsHandler(userRepository, logger),
			UpdateUser:          command.NewUpdateUserHandler(userRepository, logger),
			SetDefaultThumbnail: command.NewSetDefaultThumbnailHandler(settingRepository, logger),
			SetSetting:          command.NewSetSettingHandler(settingRepository, logger),
			AddRoleToUser:       command.NewAddRoleToUserHandler(userRepository, logger),
			AddRolesToUser:      command.NewAddRolesToUserHandler(userRepository, logger),
			RemoveRoleFromUser:  command.NewRemoveRoleFromUserHandler(userRepository, logger),
			RemoveRolesFromUser: command.NewRemoveRolesFromUserHandler(userRepository, logger),
		},
		Queries: app.Queries{
			GetTagsForMedia:        query.NewGetTagsForMediaHandler(mediaRepository, logger),
			GetAllTags:             query.NewGetAllTagsHandler(mediaRepository, logger),
			ExchangeDiscordCode:    query.NewExchangeDiscordCode(discordCtx, logger),
			GetTokenForUser:        query.NewGetTokenForUserHandler(jwtConfig, logger),
			GetDiscordUser:         query.NewGetDiscordUserHandler(discordCtx, logger),
			GetUserIdFromDiscordId: query.NewGetUserIdFromDiscordId(userRepository, logger),
			GetUserById:            query.NewGetUserById(userRepository, logger),
			GetMediaById:           query.NewGetMediaById(mediaRepository, logger),
			GetFileById:            query.NewGetFileById(fileRepository, jwtConfig, logger),
			GetMedia:               query.NewGetMedia(mediaRepository, logger),
			GetSetting:             query.NewGetSetting(settingRepository, logger),
			GetDefaultThumbnail:    query.NewGetDefaultThumbnailHandler(settingRepository, thumbnailRepository, logger),
			GetThumbnailById:       query.NewGetThumbnailById(thumbnailRepository, logger),
			GetRolesForUser:        query.NewGetRolesForUserHandler(userRepository, logger),
			GetUsers:               query.NewGetUsersHandler(userRepository, logger),
		},
		Requests: app.Requests{
			AddMedia:     request.NewAddMediaHandler(mediaRepository, userRepository, fileRepository, settingRepository, thumbnailRepository, thumbnailService, logger),
			AddThumbnail: request.NewAddThumbnailHandler(fileRepository, thumbnailRepository, thumbnailService, logger),
		},
	}
}
