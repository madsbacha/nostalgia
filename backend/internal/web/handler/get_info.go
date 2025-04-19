package handler

import (
	"encoding/json"
	"net/http"
	"nostalgia/internal/app/query"
	"nostalgia/internal/core/domain"
)

type GetInfoResponse struct {
	Title string `json:"title"`
}

func (s HttpServer) GetInfo(w http.ResponseWriter, r *http.Request) {
	title, err := s.app.Queries.GetSetting.Handle(r.Context(), query.GetSetting{
		Key: domain.SettingTitle,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to get title")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := GetInfoResponse{
		Title: title.Value,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
