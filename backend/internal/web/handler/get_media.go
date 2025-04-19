package handler

import (
	"encoding/json"
	"net/http"
	"nostalgia/internal/app/query"
	"nostalgia/internal/core/domain"
	"slices"
)

type GetMediaResponse struct {
	MediaList []Media `json:"media_list"`
}

type Media struct {
	Id                string   `json:"id"`
	Title             string   `json:"title"`
	ThumbnailUrl      string   `json:"thumbnail_url"`
	ThumbnailBlurhash string   `json:"thumbnail_blurhash"`
	Tags              []string `json:"tags"`
}

func (s HttpServer) GetMedia(w http.ResponseWriter, r *http.Request) {
	mediaList, err := s.app.Queries.GetMedia.Handle(r.Context(), query.GetMedia{})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	slices.SortStableFunc(mediaList, func(a, b domain.Media) int {
		return int(b.UploadedAt - a.UploadedAt)
	})

	responseMediaList := make([]Media, len(mediaList))
	for i, media := range mediaList {
		thumbnail, err := s.app.Queries.GetThumbnailById.Handle(r.Context(), query.GetThumbnailById{
			Id: media.ThumbnailId,
		})
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		thumbnailFile, err := s.app.Queries.GetFileById.Handle(r.Context(), query.GetFileById{
			Id: thumbnail.FileId,
		})
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		responseMediaList[i] = Media{
			Id:                media.Id,
			Title:             media.Title,
			ThumbnailUrl:      thumbnailFile.Path,
			ThumbnailBlurhash: thumbnail.BlurHash,
			Tags:              media.Tags,
		}
	}

	response := GetMediaResponse{
		MediaList: responseMediaList,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
