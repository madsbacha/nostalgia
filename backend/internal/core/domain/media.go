package domain

type Media struct {
	Id          string
	Title       string
	FileId      string
	ThumbnailId string
	UploadedAt  int64
	UploadedBy  string
	Description string
	Tags        []string
}
