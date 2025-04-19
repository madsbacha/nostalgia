package handler

import (
	"encoding/json"
	"net/http"
	"nostalgia/internal/app/query"
)

type GetAllTagsResponse struct {
	Tags []string `json:"tags"`
}

func (s HttpServer) GetAllTags(w http.ResponseWriter, r *http.Request) {
	tags, err := s.app.Queries.GetAllTags.Handle(r.Context(), query.GetAllTags{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := GetAllTagsResponse{
		Tags: tags,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
