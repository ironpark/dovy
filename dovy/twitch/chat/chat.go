package chat

import (
	"dovey/dovy/twitch/pubsub"
	"fmt"
	"github.com/gempir/go-twitch-irc/v3"
	"github.com/nicklaw5/helix/v2"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type internalEmote struct {
	img   string
	start int
	end   int
}

type ChatManager struct {
	initialized      bool
	accessToken      string
	userName         string
	displayName      string
	user             helix.User
	twitchClientId   string
	twitchApi        *helix.Client
	twitchChatClient *twitch.Client
	badges           map[string]string
	emotes           map[string]string
	bttv             []EmoteBTTV
	lock             sync.RWMutex
	msgCallback      func(message Message)
}

func NewChatManager(twitchClientId string) *ChatManager {
	emotes, _ := GetGlobalEmotes()
	cm := &ChatManager{
		initialized:      false,
		accessToken:      "",
		userName:         "",
		displayName:      "",
		user:             helix.User{},
		twitchClientId:   twitchClientId,
		twitchApi:        nil,
		twitchChatClient: nil,
		badges:           map[string]string{},
		emotes:           map[string]string{},
		bttv:             emotes,
		lock:             sync.RWMutex{},
		msgCallback:      nil,
	}

	return cm
}

func (cm *ChatManager) OnMsg(callback func(message Message)) {
	cm.lock.Lock()
	cm.msgCallback = callback
	cm.lock.Unlock()
}

func (cm *ChatManager) onMsgCallback(message twitch.PrivateMessage) {
	msg := Message{
		Id:          message.ID,
		Badges:      []string{},
		Time:        message.Time,
		Msg:         message.Message,
		UserName:    message.User.Name,
		DisplayName: message.User.DisplayName,
		Color:       message.User.Color,
	}

	for badgeId, versionId := range message.User.Badges {
		msg.Badges = append(msg.Badges, cm.getBadgeImage(badgeId, versionId))
	}
	var emotes []internalEmote
	for _, emote := range message.Emotes {
		for _, position := range emote.Positions {
			emotes = append(emotes, internalEmote{
				img:   cm.getEmoteImage(emote.ID),
				start: position.Start,
				end:   position.End,
			})
		}
	}
	if len(emotes) > 0 {
		sort.Slice(emotes, func(i, j int) bool {
			return emotes[i].start > emotes[j].start
		})
		for _, emote := range emotes {
			if emote.img == "" {
				fmt.Println("??????????????????????")
			}
			msgCopy := msg.Msg
			msg.Msg = msgCopy[:emote.start] + emote.img + msgCopy[emote.end+1:]
		}
	}
	for _, bttv := range cm.bttv {
		msg.Msg = strings.ReplaceAll(msg.Msg, bttv.Code, fmt.Sprintf("<img src=\"%s\" class=\"emote\"/>", fmt.Sprintf("https://cdn.betterttv.net/emote/%s/1x", bttv.ID)))
	}

	cm.lock.RLock()
	if cm.msgCallback != nil {
		cm.msgCallback(msg)
	}
	cm.lock.RUnlock()
}

func (cm *ChatManager) Initialize(token string) error {
	cm.lock.Lock()
	cm.accessToken = token
	defer cm.lock.Unlock()
	if cm.initialized {
		return nil
	}

	twitchApi, err := helix.NewClient(&helix.Options{
		ClientID:        cm.twitchClientId,
		ClientSecret:    "",
		AppAccessToken:  "",
		UserAccessToken: token,
		UserAgent:       "",
		RedirectURI:     "http://localhost:53324/authorize",
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
		RateLimitFunc: nil,
	})
	if err != nil {
		return err
	}
	res, err := twitchApi.GetUsers(&helix.UsersParams{
		IDs:    nil,
		Logins: nil,
	})
	if err != nil {
		return err
	}
	client := twitch.NewClient(res.Data.Users[0].Login, "oauth:"+token)
	//client.SetRateLimiter(twitch.CreateVerifiedRateLimiter())
	if err != nil {
		return err
	}
	client.OnPrivateMessage(cm.onMsgCallback)
	cm.user = res.Data.Users[0]
	cm.userName = cm.user.Login
	cm.displayName = cm.user.DisplayName
	cm.twitchApi = twitchApi
	cm.twitchChatClient = client
	// Twitch IRC chat message callback
	//client.OnNamesMessage()

	badges, _ := cm.twitchApi.GetGlobalChatBadges()
	for _, badge := range badges.Data.Badges {
		for _, version := range badge.Versions {
			cm.badges[badge.SetID+"#"+version.ID] = version.ImageUrl1x
		}
	}
	emotes, err := cm.twitchApi.GetGlobalEmotes()
	if err != nil {
		return err
	}
	for _, emote := range emotes.Data.Emotes {
		cm.emotes[emote.ID] = fmt.Sprintf("<img src=\"%s\" class=\"emote\"/>", emote.Images.Url1x)
	}
	cm.initialized = true
	go cm.twitchChatClient.Connect()
	return nil
}

func (cm *ChatManager) getEmoteImage(id string) string {
	cm.lock.RLock()
	if !cm.initialized {
		cm.lock.RUnlock()
		return ""
	}
	emote := cm.emotes[id]
	cm.lock.RUnlock()
	if emote != "" {
		return emote
	}
	cm.lock.Lock()
	cm.emotes[id] = fmt.Sprintf("<img src=\"%s\" class=\"emote\"/>", fmt.Sprintf("https://static-cdn.jtvnw.net/emoticons/v2/%s/static/light/1.0", id))
	emote = cm.emotes[id]
	cm.lock.Unlock()
	return emote
}

func (cm *ChatManager) getBadgeImage(badgeId string, versionId int) string {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	if !cm.initialized {
		return ""
	}
	return cm.badges[fmt.Sprintf("%s#%d", badgeId, versionId)]
}

func (cm *ChatManager) getChannelIdFromId(loginId string) (string, error) {
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

func (cm *ChatManager) Disconnect() {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	if cm.initialized {
		_ = cm.twitchChatClient.Disconnect()
	}
}

func (cm *ChatManager) Connect(channelName string) error {
	fmt.Println("Connect", channelName)
	cid, err := cm.getChannelIdFromId(channelName)
	if err != nil {
		return err
	}

	pubSub, err := pubsub.New()
	if err != nil {
		fmt.Println(err)
		return err
	}

	pubSub.Listen(cm.accessToken, "sjsjdjdjdjdjdjd", pubsub.TopicStreamInfo(cid))
	go pubSub.ReadLoop()

	btt, _ := GetChannelEmotesBTTV(cid)
	fmt.Println("BTTV", btt)
	cm.bttv = append(cm.bttv, btt...)
	fmt.Println("GetChannelChatBadges")
	emotes, err := cm.twitchApi.GetChannelEmotes(&helix.GetChannelEmotesParams{
		BroadcasterID: cid,
	})
	if err != nil {
		return err
	}
	for _, emote := range emotes.Data.Emotes {
		cm.emotes[emote.ID] = fmt.Sprintf("<img src=\"%s\" class=\"emote\"/>", emote.Images.Url1x)
	}

	badges, err := cm.twitchApi.GetChannelChatBadges(
		&helix.GetChatBadgeParams{
			BroadcasterID: cid,
		})
	if err != nil {
		return err
	}
	cm.lock.Lock()
	for _, badge := range badges.Data.Badges {
		for _, version := range badge.Versions {
			cm.badges[badge.SetID+"#"+version.ID] = version.ImageUrl1x
		}
	}
	cm.lock.Unlock()
	fmt.Println("tcm.lock.Unlock")
	cm.twitchChatClient.Join(channelName)
	fmt.Println("Userlist")
	userList, err := cm.twitchChatClient.Userlist(channelName)
	fmt.Println(userList)
	if err != nil {
		return err
	}

	return nil
}

//is not only chatting client
