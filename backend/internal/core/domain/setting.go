package domain

const (
	SettingKeyDefaultThumbnail = "default_thumbnail"
	SettingIsInitialized       = "initialized"
	SettingTitle               = "title"
)

type Setting struct {
	Key   string
	Value string
}
