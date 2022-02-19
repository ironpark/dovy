package emote

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type EmoteFFZ struct {
	ID   int `json:"id"`
	User struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
	} `json:"user"`
	Code   string `json:"code"`
	Images struct {
		OneX  string `json:"1x"`
		TwoX  string `json:"2x"`
		FourX string `json:"4x"`
	} `json:"images"`
	ImageType string `json:"imageType"`
}

func (efz EmoteFFZ) GetCode() string {
	return efz.Code
}

func (efz EmoteFFZ) URL1X() string {
	return efz.Images.OneX
}

func (efz EmoteFFZ) GetImgTag() string {
	return fmt.Sprintf("<img src=\"%s\" class=\"emote\" alt=\"%s\"/>", efz.URL1X(), efz.GetCode())
}

func GetChannelEmotesFFZ(channelId string) (emotes []Emote, err error) {
	fEmotes := []EmoteFFZ{}
	res, err := http.Get(fmt.Sprintf("https://api.betterttv.net/3/cached/frankerfacez/users/twitch/%s", channelId))
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(res.Body).Decode(&fEmotes)
	res.Body.Close()
	for _, emote := range fEmotes {
		emotes = append(emotes, emote)
	}
	return
}
