package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"nostalgia/internal/app/query"
)

type GetTagsForMediaResponse struct {
	Tags []string `json:"tags"`
}

func (s HttpServer) GetTagsForMedia(w http.ResponseWriter, r *http.Request) {
	mediaID := chi.URLParam(r, "mediaID")
	if mediaID == "" {
		http.Error(w, "media could not be found", http.StatusNotFound)
		return
	}

	tags, err := s.app.Queries.GetTagsForMedia.Handle(r.Context(), query.GetTagsForMedia{
		MediaId: mediaID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := GetTagsForMediaResponse{
		Tags: tags,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
