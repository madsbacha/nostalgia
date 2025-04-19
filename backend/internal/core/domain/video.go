package domain

const (
	PublicVideo = iota + 1
	UnlistedVideo
	PrivateVideo
)

type Video struct {
	Id          string
	Title       string
	UploadedAt  int64
	Visibility  int
	Views       int64
	UploadedBy  int64
	Description string
}

func (v *Video) IsPublic() bool {
	return v.Visibility == PublicVideo
}

func (v *Video) IsUnlisted() bool {
	return v.Visibility == UnlistedVideo
}

func (v *Video) IsPrivate() bool {
	return v.Visibility == PrivateVideo
}

func (v *Video) SetVisibility(visibility int) {
	v.Visibility = visibility
}
