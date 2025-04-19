package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"nostalgia/internal/app/query"
)

type GetMediaByIdResponse struct {
	Id          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Source      string   `json:"source"`
	MimeType    string   `json:"mimetype"`
	Tags        []string `json:"tags"`
}

func (s HttpServer) GetMediaById(w http.ResponseWriter, r *http.Request) {
	mediaID := chi.URLParam(r, "mediaID")
	if mediaID == "" {
		http.Error(w, "media could not be found", http.StatusNotFound)
		return
	}

	media, err := s.app.Queries.GetMediaById.Handle(r.Context(), query.GetMediaById{
		Id: mediaID,
	})
	if err != nil {
		http.Error(w, "media could not be found", http.StatusNotFound)
		return
	}

	file, err := s.app.Queries.GetFileById.Handle(r.Context(), query.GetFileById{
		Id: media.FileId,
	})
	if err != nil {
		http.Error(w, "file could not be found", http.StatusNotFound)
		return
	}

	response := GetMediaByIdResponse{
		Id:          media.Id,
		Title:       media.Title,
		Description: media.Description,
		Source:      file.Path,
		MimeType:    file.MimeType,
		Tags:        media.Tags,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
