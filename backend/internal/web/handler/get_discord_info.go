package handler

import (
	"encoding/json"
	"net/http"
	"nostalgia/internal/common/env"
)

type GetDiscordInfoResponse struct {
	ClientId     string `json:"client_id"`
	ResponseType string `json:"response_type"`
	Scopes       string `json:"scopes"`
}

func (s HttpServer) GetDiscordInfo(w http.ResponseWriter, r *http.Request) {
	res := GetDiscordInfoResponse{
		ClientId:     env.MustGet("DISCORD_CLIENT_ID"),
		ResponseType: "code",
		Scopes:       env.MustGet("DISCORD_CLIENT_SCOPES"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
