package handler

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"nostalgia/internal/app/request"
	"nostalgia/internal/common/util"
	"nostalgia/internal/web/middleware"
	"path/filepath"
)

type UploadMediaResponse struct {
	MediaId string `json:"media_id"`
}

func (s HttpServer) UploadMedia(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := middleware.UserIdFromContext(ctx)

	err := r.ParseForm()
	if err != nil {
		// TODO: Maybe not output error directly to the user
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		// TODO: Log error
		http.Error(w, "Error uploading file", http.StatusBadRequest)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			s.logger.Errorln(err)
		}
	}(file)

	extension := filepath.Ext(fileHeader.Filename)
	mimeType := fileHeader.Header.Get("Content-Type")

	title := r.Form.Get("title")
	description := r.Form.Get("description")
	tagString := r.Form.Get("tags")
	tags := util.ParseTags(tagString)

	if title == "" {
		title = fileHeader.Filename
	}

	mediaId, err := s.app.Requests.AddMedia.Handle(ctx, request.AddMedia{
		UserId:      userId,
		File:        file,
		Extension:   extension,
		MimeType:    mimeType,
		Title:       title,
		Description: description,
		Tags:        tags,
	})
	if err != nil {
		// TODO: Do not output error to user
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := UploadMediaResponse{
		MediaId: mediaId,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
