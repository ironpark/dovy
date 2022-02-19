package twitch

import (
	"dovey/dovy/twitch/badge"
	"dovey/dovy/twitch/emote"
	"dovey/dovy/twitch/pubsub"
	odered "dovey/pkg/odered"
	"fmt"
	twitchIrc "github.com/gempir/go-twitch-irc/v3"
	"github.com/nicklaw5/helix/v2"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ConnectionManager struct {
	channels          *odered.OrderedMap[string, *Channel]
	ps                *pubsub.PubSub
	twitchApi         *helix.Client
	lock              sync.RWMutex
	token             string
	irc               *twitchIrc.Client
	bttv              *emote.Store
	callback          func(message Message)
	userEventCallback func(channel *Channel, event string, user string)
	globalBadges      *badge.Store
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
	bttvEmotes, err := emote.GetGlobalEmotesBTTV()
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
		bttv:         emote.NewStore(),
		callback:     nil,
		globalBadges: badge.NewStore(),
	}
	cm.bttv.SetEmotes(bttvEmotes)
	return cm, nil
}

func (cm *ConnectionManager) msgParse(id, channelName, message string, user twitchIrc.User, emotes []*twitchIrc.Emote, MsgTime time.Time) (msg Message, err error) {
	channel, ok := cm.channels.Get(channelName)
	if !ok {
		return Message{}, fmt.Errorf("channel not exist")
	}
	parserdMsg := message
	parts := strings.Split(message, " ")

	for i, part := range parts {
		founded := false
		if len(emotes) > 0 {
			for _, emo := range emotes {
				if part == emo.Name {
					imgUri := fmt.Sprintf("https://static-cdn.jtvnw.net/emoticons/v2/%s/default/light/1.0", emo.ID)
					parts[i] = emote.ImgTag(imgUri, emo.Name)
					founded = true
					break
				}
			}
		}
		if founded {
			continue
		}
		// TWITCH Emotes
		for _, emo := range emotes {
			if part == emo.Name {
				founded = true
				imgUri := fmt.Sprintf("https://static-cdn.jtvnw.net/emoticons/v2/%s/default/light/1.0", emo.ID)
				parts[i] = emote.ImgTag(imgUri, emo.Name)
				break
			}
		}
		if founded {
			continue
		}

		// BTTV Global Emotes
		if bttvEmote := cm.bttv.GetEmote(part); bttvEmote != nil {
			parts[i] = bttvEmote.GetImgTag()
			continue
		}

		// BTTV, FFZ Channel Emotes
		if channelEmote := channel.Emotes.GetEmote(part); channelEmote != nil {
			parts[i] = channelEmote.GetImgTag()
		}
	}
	parserdMsg = strings.Join(parts, " ")
	msg = Message{
		Id:             id,
		Channel:        channelName,
		Time:           MsgTime,
		Badges:         nil,
		Msg:            message,
		MsgWithEmotes:  parserdMsg,
		UserName:       user.Name,
		DisplayName:    user.DisplayName,
		IsUserNameOnly: strings.ToLower(user.DisplayName) == user.Name,
		Color:          user.Color,
	}

	for badgeId, versionId := range user.Badges {
		imgUri := cm.globalBadges.GetBadgeImage(badgeId, versionId)
		if imgUri == "" {
			imgUri = channel.Badges.GetBadgeImage(badgeId, versionId)
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
	fmt.Println("oauth:" + token)
	client := twitchIrc.NewClient(res.Data.Users[0].Login, "oauth:"+token)
	client.OnPrivateMessage(func(iMsg twitchIrc.PrivateMessage) {
		channel, err := cm.ConnectOrGet(iMsg.Channel)
		if err != nil {
			log.Println("ERR", err)
			return
		}

		if channel.AddUser(iMsg.User.Name) {
			fmt.Println(iMsg.User.Name, "JOIN")
			cm.lock.RLock()
			if cm.userEventCallback != nil {
				cm.userEventCallback(channel, "JOIN", iMsg.User.Name)
			}
			cm.lock.RUnlock()
		}
		msg, err := cm.msgParse(iMsg.ID, iMsg.Channel, iMsg.Message, iMsg.User, iMsg.Emotes, iMsg.Time)
		if err == nil {
			cm.lock.RLock()
			if cm.callback != nil {
				cm.callback(msg)
			}
			cm.lock.RUnlock()
		}
	})
	client.OnNamesMessage(func(iMsg twitchIrc.NamesMessage) {
		fmt.Println("[OnNames]", iMsg.Raw)
	})
	//Join a channel.
	client.OnUserJoinMessage(func(iMsg twitchIrc.UserJoinMessage) {
		channel, err := cm.ConnectOrGet(iMsg.Channel)
		if err != nil {
			log.Println("ERR", err)
			return
		}
		if channel.AddUser(iMsg.User) {
			cm.lock.RLock()
			if cm.userEventCallback != nil {
				cm.userEventCallback(channel, "JOIN", iMsg.User)
			}
			cm.lock.RUnlock()
		}
	})
	//Leave a channel.
	client.OnUserPartMessage(func(iMsg twitchIrc.UserPartMessage) {
		channel, err := cm.ConnectOrGet(iMsg.Channel)
		if err != nil {
			log.Println("ERR", err)
			return
		}
		if channel.RemoveUser(iMsg.User) {
			cm.lock.RLock()
			if cm.userEventCallback != nil {
				cm.userEventCallback(channel, "LEAVE", iMsg.User)
			}
			cm.lock.RUnlock()
		}
	})

	client.OnWhisperMessage(func(iMsg twitchIrc.WhisperMessage) {
		fmt.Println("[Whisper]", iMsg.User, iMsg.Message)
	})
	client.OnUserStateMessage(func(iMsg twitchIrc.UserStateMessage) {
		fmt.Println("[UserState]", iMsg.User.Name, iMsg.Message, iMsg.Raw)
	})
	client.OnUserNoticeMessage(func(iMsg twitchIrc.UserNoticeMessage) {
		fmt.Println("[UserNotice]", iMsg.User.Name, iMsg.Message, iMsg.Raw)
		cm.msgParse(iMsg.ID, iMsg.Channel, iMsg.Message, iMsg.User, iMsg.Emotes, iMsg.Time)
	})

	client.OnUnsetMessage(func(message twitchIrc.RawMessage) {
		fmt.Println("[OnUnsetMessage]", message.Raw)
	})
	client.Capabilities = []string{twitchIrc.CommandsCapability, twitchIrc.TagsCapability, twitchIrc.MembershipCapability}
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

func (cm *ConnectionManager) OnUserEvent(callback func(channel *Channel, event string, user string)) {
	cm.lock.Lock()
	cm.userEventCallback = callback
	cm.lock.Unlock()
}

func (cm *ConnectionManager) ConnectOrGet(channelName string) (*Channel, error) {
	v, ok := cm.channels.Get(channelName)
	if ok {
		return v, nil
	}
	channel, err := NewChannel(cm.twitchApi, channelName)
	if err != nil {
		return nil, err
	}

	cm.channels.Set(channelName, channel)
	cm.irc.Join(channelName)
	return channel, nil
}

func (cm *ConnectionManager) IsInitialized() bool {
	fmt.Println("IsInitialized")
	cm.lock.RLock()
	cm.lock.RUnlock()
	return cm.token != ""
}

func (cm *ConnectionManager) Send(channel string, msg string) {
	cm.irc.Say(channel, msg)
}
