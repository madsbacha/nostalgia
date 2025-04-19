package handler

import (
	"encoding/json"
	"net/http"
	"nostalgia/internal/app/command"
	"nostalgia/internal/app/query"
)

type AuthDiscordRequest struct {
	Code        string `json:"code"`
	RedirectUri string `json:"redirect_uri"`
}

type AuthDiscordResponse struct {
	Token string `json:"token"`
}

func (s HttpServer) AuthDiscord(w http.ResponseWriter, r *http.Request) {
	// TODO: Refactor errors

	// Parse the JSON body into AuthDiscordRequest
	var req AuthDiscordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Check if the code is provided
	if req.Code == "" {
		http.Error(w, "Missing code in request", http.StatusBadRequest)
		return
	}

	token, err := s.app.Queries.ExchangeDiscordCode.Handle(r.Context(), query.ExchangeDiscordCode{
		Code:        req.Code,
		RedirectUri: req.RedirectUri,
	})
	if err != nil {
		http.Error(w, "Failed to exchange discord code", http.StatusBadRequest)
		return
	}

	discord, err := s.app.Queries.GetDiscordUser.Handle(r.Context(), query.GetDiscordUser{
		Token: *token,
	})
	if err != nil {
		http.Error(w, "Failed to get discord user", http.StatusBadRequest)
		return
	}

	err = s.app.Commands.EnsureUserExists.Handle(r.Context(), command.EnsureUserExists{
		DiscordId: discord.Id,
		Username:  discord.Username,
	})
	if err != nil {
		http.Error(w, "Failed to ensure user exists", http.StatusBadRequest)
		return
	}

	userId, err := s.app.Queries.GetUserIdFromDiscordId.Handle(r.Context(), query.GetUserIdFromDiscordId{
		DiscordId: discord.Id,
	})
	if err != nil {
		http.Error(w, "Failed to get user id", http.StatusBadRequest)
		return
	}

	err = s.app.Commands.UpdateUser.Handle(r.Context(), command.UpdateUser{
		Id:     userId,
		Avatar: discord.Avatar,
	})
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusBadRequest)
		return
	}

	jwtToken, err := s.app.Queries.GetTokenForUser.Handle(r.Context(), query.GetTokenForUser{
		UserId: userId,
	})
	if err != nil {
		http.Error(w, "Failed to get token for user", http.StatusBadRequest)
		return
	}

	response := AuthDiscordResponse{
		Token: jwtToken,
	}

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
