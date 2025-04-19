package actions

import (
	"context"
	"database/sql"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"nostalgia/internal/cli/mediaviewer/datab"
	"nostalgia/internal/common/rbac"
	"nostalgia/internal/common/util"
	"nostalgia/internal/core/port"
	"nostalgia/internal/infrastructure"
	"nostalgia/pkg/dataquery"
	"os"
	"path/filepath"
	"strings"
)

type MigrateMediaviewerContext struct {
	context.Context
	Logger           *logrus.Entry
	OldDb            *datab.Context
	NewDb            *sql.DB
	OldStoragePath   string
	FileRepo         port.FileRepository
	UserRepo         port.UserRepository
	ThumbnailService port.ThumbnailService
	ThumbnailRepo    port.ThumbnailRepository
	MediaRepo        port.MediaRepository
	SettingRepo      port.SettingRepository
}

func MigrateMediaviewer(ctx *cli.Context) error {
	databasePath := ctx.String("database")
	storagePath := ctx.String("storage")
	ffmpegBinaryPath := ctx.String("ffmpeg-binary")
	oldDatabasePath := ctx.String("old-database")
	oldStoragePath := ctx.String("old-storage")

	tempDirPath := ctx.String("temporary-directory")

	logger := logrus.NewEntry(logrus.StandardLogger())

	logger.Printf("%s -> %s\n", oldDatabasePath, databasePath)

	oldDb := datab.New(oldDatabasePath)

	db := infrastructure.NewSqliteDatabase(databasePath)
	settingRepository := infrastructure.NewSettingSqliteRepository(db)
	mediaRepository := infrastructure.NewMediaSqliteRepository(db)
	userRepository := infrastructure.NewUserSqliteRepository(db)
	thumbnailRepository := infrastructure.NewThumbnailSqliteRepository(db)
	fileRepository := infrastructure.NewFileLocalRepository(db, storagePath)
	thumbnailService := infrastructure.NewThumbnailService(ffmpegBinaryPath, tempDirPath, logger)

	migrateCtx := MigrateMediaviewerContext{
		Context:          ctx.Context,
		Logger:           logger,
		OldDb:            &oldDb,
		NewDb:            db,
		OldStoragePath:   oldStoragePath,
		FileRepo:         fileRepository,
		UserRepo:         userRepository,
		ThumbnailService: thumbnailService,
		ThumbnailRepo:    thumbnailRepository,
		MediaRepo:        mediaRepository,
		SettingRepo:      settingRepository,
	}

	migrateUsers(&migrateCtx)
	migrateMedia(&migrateCtx)

	return nil
}

func migrateUsers(ctx *MigrateMediaviewerContext) {
	res, err := ctx.UserRepo.GetAll(ctx)
	if (err == nil && len(res) > 0) || (err != nil && !errors.Is(err, sql.ErrNoRows)) {
		ctx.Logger.Fatalln("database contains existing users. Please start from a clean slate.")
	}

	ctx.Logger.Println("migrating users")
	oldUsers, err := datab.Query[*datab.User](&ctx.OldDb.User).GetAll()
	if err != nil {
		ctx.Logger.WithError(err).Fatalln("Failed to query users")
	}

	for _, oldUser := range oldUsers {
		var roles []string
		if oldUser.Id == 1 {
			roles = append(roles, rbac.RoleWhitelisted)
			roles = append(roles, rbac.RoleCanManagePermissions)
		}
		err := ctx.UserRepo.Insert(ctx, port.NewUser{
			DiscordId: oldUser.DiscordId,
			Username:  oldUser.Username,
			Avatar:    oldUser.Avatar,
			Roles:     roles,
		})
		if err != nil {
			ctx.Logger.WithError(err).Fatalf("Failed to insert user %s", oldUser.DiscordId)
		}
	}
	ctx.Logger.Println("migrated users")
}

func migrateMedia(ctx *MigrateMediaviewerContext) {
	ctx.Logger.Println("migrating media")

	oldMedia, err := datab.Query[*datab.Media](&ctx.OldDb.Media).GetAll()
	if err != nil {
		ctx.Logger.WithError(err).Fatalln("Failed to query media")
	}

	progress := 0
	maxProgress := len(oldMedia)
	scanner := infrastructure.NewMediaScanner()
	for _, oldMediaItem := range oldMedia {
		if oldMediaItem.Deleted {
			ctx.Logger.Infof("Skipping deleted media %s\n", oldMediaItem.Id)
			maxProgress--
			continue
		}

		objId := oldMediaItem.ObjectId
		bucketObject, err := datab.Query[*datab.BucketObject](&ctx.OldDb.Bucket).Where("id = ?", objId).GetFirst()
		if err != nil {
			ctx.Logger.WithError(err).Fatalf("Failed to query for bucket %d\n", objId)
		}
		filename := bucketObject.Name

		if !isVideoExtension(filepath.Ext(filename)) {
			ctx.Logger.Infof("Skipping non-video file: %s\n", filename)
			maxProgress--
			continue
		}

		filePath := filepath.Join(ctx.OldStoragePath, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			ctx.Logger.Errorf("File %s does not exist, skipping.", filePath)
			maxProgress--
			continue
		}

		oldFile, err := os.OpenFile(filePath, os.O_RDWR, 0644)
		if err != nil {
			ctx.Logger.WithError(err).Fatalf("Failed to open file %s", filePath)
		}

		mimeType := getMimeTypeFromExtension(filepath.Ext(filename))
		newFileId, err := ctx.FileRepo.Upload(ctx.Context, oldFile, mimeType, filepath.Ext(filename))
		if err != nil {
			ctx.Logger.WithError(err).Fatalf("Failed to upload file %s", filePath)
		}

		generatedThumbnail, err := ctx.ThumbnailService.GenerateFromVideo(ctx.Context, filePath)
		if err != nil {
			ctx.Logger.WithError(err).Fatalf("Failed to generate thumbnail for file %s", filePath)
		}

		newThumbnailFileId, err := ctx.FileRepo.Upload(ctx.Context, generatedThumbnail.File, generatedThumbnail.MimeType, generatedThumbnail.Extension)
		if err != nil {
			ctx.Logger.WithError(err).Fatalf("Failed to upload thumbnail for file %s", filePath)
		}

		thumbnailId, err := ctx.ThumbnailRepo.Create(ctx.Context, port.NewThumbnail{
			FileId:   newThumbnailFileId,
			Blurhash: generatedThumbnail.Blurhash,
		})
		if err != nil {
			ctx.Logger.WithError(err).Fatalf("Failed to create thumbnail for file %s", filePath)
		}

		mediaQuery := dataquery.QueryContext[*infrastructure.Media](ctx, ctx.NewDb, infrastructure.TableMedia, scanner)
		uploadedBy := oldMediaItem.UploadedBy
		if uploadedBy == 0 {
			uploadedBy = 1
		}
		_, err = mediaQuery.Insert(map[string]interface{}{
			infrastructure.MediaColumnId:          oldMediaItem.Id,
			infrastructure.MediaColumnTitle:       oldMediaItem.Title,
			infrastructure.MediaColumnDescription: oldMediaItem.Description,
			infrastructure.MediaColumnFileId:      newFileId,
			infrastructure.MediaColumnUploadedBy:  uploadedBy,
			infrastructure.MediaColumnThumbnailId: thumbnailId,
			infrastructure.MediaColumnTags:        util.JoinTags(oldMediaItem.Tags),
			infrastructure.MediaColumnUploadedAt:  oldMediaItem.UploadedAt,
		})
		if err != nil {
			ctx.Logger.WithError(err).Fatalf("Failed to insert media %s", oldMediaItem.Id)
		}

		ctx.Logger.Printf("Migrated %d/%d media\n", progress, maxProgress)
		progress++
	}

	ctx.Logger.Println("migrated media")
}

func isVideoExtension(extension string) bool {
	extension = strings.ToLower(extension)
	switch extension {
	case ".mp4", ".mov", ".mkv", ".webm":
		return true
	default:
		return false
	}
}

func getMimeTypeFromExtension(extension string) string {
	extension = strings.ToLower(extension)
	switch extension {
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".mp3":
		return "audio/mpeg"
	case ".mov":
		return "video/quicktime"
	case ".mkv":
		return "video/x-matroska"
	default:
		return "application/octet-stream"
	}
}
