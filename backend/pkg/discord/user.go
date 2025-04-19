package discord

import "encoding/json"

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

func UserFromJson(s []byte) (User, error) {
	user := User{}

	err := json.Unmarshal(s, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
