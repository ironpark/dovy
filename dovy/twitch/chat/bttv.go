package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type EmoteBTTV struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	ImageType string `json:"imageType"`
	UserID    string `json:"userId"`
}

//https://api.betterttv.net/3/cached/users/twitch/121059319
func GetGlobalEmotes() (emotes []EmoteBTTV, err error) {
	res, err := http.Get("https://api.betterttv.net/3/cached/emotes/global")
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(res.Body).Decode(&emotes)
	res.Body.Close()
	return
}

func GetChannelEmotesBTTV(channelId string) (emotes []EmoteBTTV, err error) {
	res, err := http.Get(fmt.Sprintf("https://api.betterttv.net/3/cached/users/twitch/%s", channelId))
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(res.Body).Decode(&emotes)
	res.Body.Close()

	return
}

func GetChannelEmotesFFZ(channelId string) (emotes []EmoteBTTV, err error) {
	res, err := http.Get(fmt.Sprintf("https://api.betterttv.net/3/cached/frankerfacez/users/twitch/%s", channelId))
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(res.Body).Decode(&emotes)
	res.Body.Close()

	return
}
