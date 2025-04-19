package discord

import (
	"fmt"
)

func AvatarUrl(discordId string, avatarHash string) string {
	if avatarHash == "" {
		return ""
	}
	return fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", discordId, avatarHash)
}
