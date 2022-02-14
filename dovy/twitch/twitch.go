package twitch

import (
	"dovey/dovy/twitch/badge"
	"dovey/dovy/twitch/emote"
	"dovey/dovy/twitch/pubsub"
	odered "dovey/pkg/odered"
	"fmt"
	twitchIrc "github.com/gempir/go-twitch-irc/v3"
	"github.com/nicklaw5/helix/v2"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ConnectionManager struct {
	channels     *odered.OrderedMap[string, *Channel]
	ps           *pubsub.PubSub
	twitchApi    *helix.Client
	lock         sync.RWMutex
	token        string
	irc          *twitchIrc.Client
	bttv         []emote.EmoteBTTV
	callback     func(message Message)
	globalBadges *badge.Store
}

func NewConnectionManager() (*ConnectionManager, error) {
	ps, err := pubsub.New()
	if err != nil {
		return nil, err
	}
	twitchApi, err := helix.NewClient(&helix.Options{
		ClientID:       TwitchClientId,
		ClientSecret:   "",
		AppAccessToken: "",
		UserAgent:      "",
		RedirectURI:    "http://localhost:53324/authorize",
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
		RateLimitFunc: nil,
	})
	if err != nil {
		return nil, err
	}

	cm := &ConnectionManager{
		channels:     odered.NewMap[string, *Channel](),
		ps:           ps,
		twitchApi:    twitchApi,
		lock:         sync.RWMutex{},
		token:        "",
		irc:          nil,
		bttv:         nil,
		callback:     nil,
		globalBadges: badge.New(),
	}
	return cm, nil
}
func (cm *ConnectionManager) msgParse(message twitchIrc.PrivateMessage) (msg Message, err error) {
	channel, ok := cm.channels.Get(message.Channel)
	if !ok {
		return Message{}, fmt.Errorf("channel not exist")
	}
	parserdMsg := message.Message
	parts := strings.Split(message.Message, " ")
	if len(message.Emotes) > 0 {
		for _, emo := range message.Emotes {
			for i, part := range parts {
				if part == emo.Name {
					imgUri := fmt.Sprintf("https://static-cdn.jtvnw.net/emoticons/v2/%s/static/light/1.0", emo.ID)
					parts[i] = emote.ImgTag(imgUri)
				}
			}
		}
	}

	for _, bttv := range cm.bttv {
		for i, part := range parts {
			if part == bttv.Code {
				imgUri := fmt.Sprintf("https://cdn.betterttv.net/emote/%s/1x", bttv.ID)
				parts[i] = emote.ImgTag(imgUri)
			}
		}
	}

	parserdMsg = strings.Join(parts, " ")
	msg = Message{
		Id:             message.ID,
		Channel:        message.Channel,
		Time:           time.Time{},
		Badges:         nil,
		Msg:            message.Message,
		MsgWithEmotes:  parserdMsg,
		UserName:       message.User.Name,
		DisplayName:    message.User.DisplayName,
		IsUserNameOnly: strings.ToLower(message.User.DisplayName) == message.User.Name,
		Color:          message.User.Color,
	}

	for badgeId, versionId := range message.User.Badges {
		imgUri := cm.globalBadges.GetBadgeImage(badgeId, versionId)
		if imgUri == "" {
			imgUri = channel.badgeStore.GetBadgeImage(badgeId, versionId)
		}
		if imgUri == "" {
			continue
		}
		msg.Badges = append(msg.Badges, imgUri)
	}
	return msg, nil
}
func (cm *ConnectionManager) Initialize(token string) error {
	// initialize check
	cm.lock.RLock()
	if cm.token != "" {
		cm.lock.RUnlock()
		return nil
	}
	cm.lock.RUnlock()
	// set twitch user token (helix api)
	cm.twitchApi.SetUserAccessToken(token)

	res, err := cm.twitchApi.GetUsers(&helix.UsersParams{})
	if err != nil {
		return err
	}
	//set global badges
	globalBadge, err := cm.twitchApi.GetGlobalChatBadges()
	if err != nil {
		return err
	}
	cm.globalBadges.SetBadges(globalBadge.Data.Badges)

	client := twitchIrc.NewClient(res.Data.Users[0].Login, "oauth:"+token)
	client.OnPrivateMessage(func(message twitchIrc.PrivateMessage) {
		msg, err := cm.msgParse(message)
		if err == nil {
			cm.lock.RLock()
			if cm.callback != nil {
				cm.callback(msg)
			}
			cm.lock.RUnlock()
		}
	})

	client.OnUserJoinMessage(func(message twitchIrc.UserJoinMessage) {
		fmt.Println("[UserJoin]", message.Channel, message.User)
	})

	client.OnWhisperMessage(func(message twitchIrc.WhisperMessage) {
		fmt.Println("[Whisper]", message.User, message.Message)
	})

	go func() {
		err = client.Connect()
		if err != nil {
			fmt.Println(err)
		}
	}()

	cm.lock.Lock()
	cm.irc = client
	cm.token = token
	cm.lock.Unlock()
	return nil
}

func (cm *ConnectionManager) OnMessage(callback func(message Message)) {
	cm.lock.Lock()
	cm.callback = callback
	cm.lock.Unlock()
}

func (cm *ConnectionManager) getChannelIdFromId(loginId string) (string, error) {
	users, err := cm.twitchApi.GetUsers(&helix.UsersParams{
		IDs:    nil,
		Logins: []string{loginId},
	})
	if err != nil {
		return "", err
	}
	if users.Error != "" {
		return "", fmt.Errorf(users.ErrorMessage)
	}
	if len(users.Data.Users) == 0 {
		return "", fmt.Errorf("Not Exists")
	}
	return users.Data.Users[0].ID, nil
}

func (cm *ConnectionManager) ConnectOrGet(channelName string) (*Channel, error) {
	fmt.Println("ConnectOrGet", channelName)
	v, ok := cm.channels.Get(channelName)
	if ok {
		return v, nil
	}
	channelId, err := cm.getChannelIdFromId(channelName)
	if err != nil {
		return nil, err
	}
	channel := &Channel{
		RWMutex:    sync.RWMutex{},
		streamer:   channelName,
		channelId:  channelId,
		viewers:    0,
		users:      odered.NewSet[string](),
		emoteStore: nil,
		badgeStore: badge.New(),
	}
	badges, err := cm.twitchApi.GetChannelChatBadges(
		&helix.GetChatBadgeParams{
			BroadcasterID: channelId,
		})
	if err != nil {
		return nil, err
	}
	channel.badgeStore.SetBadges(badges.Data.Badges)
	cm.channels.Set(channelName, channel)
	cm.irc.Join(channelName)
	return channel, nil
}
