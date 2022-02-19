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
	bttv         *emote.Store
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
	if len(emotes) > 0 {
		for _, emo := range emotes {
			for i, part := range parts {
				if part == emo.Name {
					imgUri := fmt.Sprintf("https://static-cdn.jtvnw.net/emoticons/v2/%s/default/light/1.0", emo.ID)
					parts[i] = emote.ImgTag(imgUri)
				}
			}
		}
	}
	for i, part := range parts {
		founded := false
		// TWITCH Emotes
		for _, emo := range emotes {
			if part == emo.Name {
				founded = true
				imgUri := fmt.Sprintf("https://static-cdn.jtvnw.net/emoticons/v2/%s/default/light/1.0", emo.ID)
				parts[i] = emote.ImgTag(imgUri)
				break
			}
		}
		if founded {
			continue
		}
		// BTTV Global Emotes
		if emoteUrl := cm.bttv.GetEmote(part); emoteUrl != "" {
			parts[i] = emote.ImgTag(emoteUrl)
			continue
		}
		// BTTV, FFZ Channel Emotes
		if emoteUrl := channel.emoteStore.GetEmote(part); emoteUrl != "" {
			parts[i] = emote.ImgTag(emoteUrl)
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
	fmt.Println("oauth:" + token)
	client := twitchIrc.NewClient(res.Data.Users[0].Login, "oauth:"+token)
	client.OnPrivateMessage(func(iMsg twitchIrc.PrivateMessage) {
		msg, err := cm.msgParse(iMsg.ID, iMsg.Channel, iMsg.Message, iMsg.User, iMsg.Emotes, iMsg.Time)
		if err == nil {
			cm.lock.RLock()
			if cm.callback != nil {
				cm.callback(msg)
			}
			cm.lock.RUnlock()
		}
	})
	client.OnNamesMessage(func(message twitchIrc.NamesMessage) {
		fmt.Println("[OnNames]", message.Raw)
	})
	client.OnUserJoinMessage(func(message twitchIrc.UserJoinMessage) {
		fmt.Println("[UserJoin]", message.Channel, message.User)
	})
	client.OnWhisperMessage(func(message twitchIrc.WhisperMessage) {
		fmt.Println("[Whisper]", message.User, message.Message)
	})
	client.OnUserStateMessage(func(message twitchIrc.UserStateMessage) {
		fmt.Println("[UserState]", message.User.Name, message.Message, message.Raw)
	})
	client.OnUserNoticeMessage(func(iMsg twitchIrc.UserNoticeMessage) {
		fmt.Println("[UserNotice]", iMsg.User.Name, iMsg.Message, iMsg.Raw)
		cm.msgParse(iMsg.ID, iMsg.Channel, iMsg.Message, iMsg.User, iMsg.Emotes, iMsg.Time)
	})
	client.OnUserPartMessage(func(message twitchIrc.UserPartMessage) {
		fmt.Println("[OnUserPartMessage]", message.Raw)
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
		return "", fmt.Errorf("not exists")
	}
	return users.Data.Users[0].ID, nil
}

func (cm *ConnectionManager) ConnectOrGet(channelName string) (*Channel, error) {
	v, ok := cm.channels.Get(channelName)
	if ok {
		return v, nil
	}
	channelId, err := cm.getChannelIdFromId(channelName)
	if err != nil {
		fmt.Println("ERR!", err)
		return nil, err
	}

	channel := &Channel{
		RWMutex:    sync.RWMutex{},
		streamer:   channelName,
		channelId:  channelId,
		viewers:    0,
		users:      odered.NewSet[string](),
		emoteStore: emote.NewStore(),
		badgeStore: badge.NewStore(),
	}
	badges, err := cm.twitchApi.GetChannelChatBadges(
		&helix.GetChatBadgeParams{
			BroadcasterID: channelId,
		})
	if err != nil {
		fmt.Println("ERR!", err)
		return nil, err
	}

	bttvChannelEmotes, err := emote.GetChannelEmotesBTTV(channelId)
	if err != nil {
		fmt.Println("ERR!", err)
		return nil, err
	}
	ffzChannelEmotes, err := emote.GetChannelEmotesFFZ(channelId)
	if err != nil {
		fmt.Println("ERR!", err)
		return nil, err
	}
	channel.emoteStore.SetEmotes(bttvChannelEmotes)
	channel.emoteStore.SetEmotes(ffzChannelEmotes)
	channel.badgeStore.SetBadges(badges.Data.Badges)
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
