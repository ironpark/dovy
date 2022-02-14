package twitch

import "time"

type Message struct {
	Id             string    `json:"id"`
	Channel        string    `json:"channel"`
	Time           time.Time `json:"time"`
	Badges         []string  `json:"badges"`
	Msg            string    `json:"msg"`
	MsgWithEmotes  string    `json:"msg_with_emotes"`
	UserName       string    `json:"user_name"`
	DisplayName    string    `json:"display_name"`
	IsUserNameOnly bool      `json:"is_user_name_only"`
	Color          string    `json:"color"`
}
