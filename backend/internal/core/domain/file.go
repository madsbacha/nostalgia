package domain

type File struct {
	Id           string `json:"id"`
	MimeType     string `json:"mime_type"`
	Path         string `json:"path"`
	InternalPath string `json:"internal_path"`
}
