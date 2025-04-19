package handler

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"nostalgia/internal/app/query"
	"nostalgia/internal/web/middleware"
)

func (s HttpServer) GetFileById(w http.ResponseWriter, r *http.Request) {
	fileIdFromToken := middleware.FileIdFromContext(r.Context())
	fileId := chi.URLParam(r, "fileID")

	if fileId != "" && fileId != fileIdFromToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if fileId == "" {
		http.Error(w, "file could not be found", http.StatusNotFound)
		return
	}

	file, err := s.app.Queries.GetFileById.Handle(r.Context(), query.GetFileById{
		Id: fileId,
	})
	if err != nil {
		http.Error(w, "file could not be found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, file.InternalPath)
}
