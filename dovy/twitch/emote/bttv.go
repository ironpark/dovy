package emote

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type EmoteBTTV struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	ImageType string `json:"imageType"`
}

type channelEmoteResp struct {
	ID            string        `json:"id"`
	Bots          []interface{} `json:"bots"`
	Avatar        string        `json:"avatar"`
	ChannelEmotes []EmoteBTTV   `json:"channelEmotes"`
	SharedEmotes  []EmoteBTTV   `json:"sharedEmotes"`
}

func (ebttv EmoteBTTV) GetCode() string {
	return ebttv.Code
}

func (ebttv EmoteBTTV) URL1X() string {
	return fmt.Sprintf("https://cdn.betterttv.net/emote/%s/1x", ebttv.ID)
}

func GetGlobalEmotesBTTV() (emotes []Emote, err error) {
	bEmotes := []EmoteBTTV{}
	res, err := http.Get("https://api.betterttv.net/3/cached/emotes/global")
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(res.Body).Decode(&bEmotes)
	res.Body.Close()
	for _, emote := range bEmotes {
		emotes = append(emotes, Emote(emote))
	}
	return emotes, nil
}

func GetChannelEmotesBTTV(channelId string) (emotes []Emote, err error) {
	var bEmotes channelEmoteResp
	fmt.Println("GetChannelEmotesBTTV", fmt.Sprintf("https://api.betterttv.net/3/cached/users/twitch/%s", channelId))
	res, err := http.Get(fmt.Sprintf("https://api.betterttv.net/3/cached/users/twitch/%s", channelId))
	if err != nil {
		fmt.Println("GetChannelEmotesBTTV", err)
		return nil, err
	}
	err = json.NewDecoder(res.Body).Decode(&bEmotes)
	res.Body.Close()
	if err != nil {
		fmt.Println("GetChannelEmotesBTTV", err)
		return nil, err
	}
	for _, emote := range bEmotes.ChannelEmotes {
		emotes = append(emotes, Emote(emote))
	}
	for _, emote := range bEmotes.SharedEmotes {
		emotes = append(emotes, Emote(emote))
	}
	return
}
