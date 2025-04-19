package app

import (
	"nostalgia/internal/app/command"
	"nostalgia/internal/app/query"
	"nostalgia/internal/app/request"
)

type Application struct {
	Commands Commands
	Queries  Queries
	Requests Requests
}

type Commands struct {
	AddTagToMedia       command.AddTagToMediaHandler
	RemoveTagFromMedia  command.RemoveTagFromMediaHandler
	EnsureUserExists    command.EnsureUserExistsHandler
	UpdateUser          command.UpdateUserHandler
	SetDefaultThumbnail command.SetDefaultThumbnailHandler
	SetSetting          command.SetSettingHandler
	AddRoleToUser       command.AddRoleToUserHandler
	AddRolesToUser      command.AddRolesToUserHandler
	RemoveRoleFromUser  command.RemoveRoleFromUserHandler
	RemoveRolesFromUser command.RemoveRolesFromUserHandler
}

type Queries struct {
	GetTagsForMedia        query.GetTagsForMediaHandler
	GetAllTags             query.GetAllTagsHandler
	ExchangeDiscordCode    query.ExchangeDiscordCodeHandler
	GetTokenForUser        query.GetTokenForUserHandler
	GetDiscordUser         query.GetDiscordUserHandler
	GetUserIdFromDiscordId query.GetUserIdFromDiscordIdHandler
	GetUserById            query.GetUserByIdHandler
	GetMediaById           query.GetMediaByIdHandler
	GetFileById            query.GetFileByIdHandler
	GetMedia               query.GetMediaHandler
	GetSetting             query.GetSettingHandler
	GetDefaultThumbnail    query.GetDefaultThumbnailHandler
	GetThumbnailById       query.GetThumbnailByIdHandler
	GetRolesForUser        query.GetRolesForUserHandler
	GetUsers               query.GetUsersHandler
}

type Requests struct {
	AddMedia     request.AddMediaHandler
	AddThumbnail request.AddThumbnailHandler
}
