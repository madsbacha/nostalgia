package handler

import (
	"encoding/json"
	"net/http"
	"nostalgia/internal/app/query"
	"nostalgia/internal/web/middleware"
	"nostalgia/internal/web/models"
	"nostalgia/pkg/discord"
)

type CurrentUserResponse struct {
	Id          string                 `json:"id"`
	Username    string                 `json:"username"`
	AvatarUrl   string                 `json:"avatar_url"`
	Permissions models.UserPermissions `json:"permissions"`
}

func (s HttpServer) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := middleware.UserIdFromContext(ctx)
	user, err := s.app.Queries.GetUserById.Handle(ctx, query.GetUserById{
		Id: userId,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := CurrentUserResponse{
		Id:          user.Id,
		Username:    user.Username,
		AvatarUrl:   discord.AvatarUrl(user.DiscordId, user.Avatar),
		Permissions: models.NewUserPermissionsFromRoles(user.Roles),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
