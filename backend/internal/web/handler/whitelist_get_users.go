package handler

import (
	"encoding/json"
	"net/http"
	"nostalgia/internal/app/query"
	"nostalgia/pkg/discord"
)

type GetUsersResponse struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
}

func (s HttpServer) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := s.app.Queries.GetUsers.Handle(ctx, query.GetUsers{})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := make([]GetUsersResponse, len(users))
	for i, user := range users {
		response[i] = GetUsersResponse{
			Id:        user.Id,
			Username:  user.Username,
			AvatarUrl: discord.AvatarUrl(user.DiscordId, user.Avatar),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
