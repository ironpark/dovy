package twitch

import (
	"dovey/dovy/twitch/badge"
	"dovey/dovy/twitch/emote"
	odered "dovey/pkg/odered"
	"fmt"
	"github.com/nicklaw5/helix/v2"
	"sync"
)

type Channel struct {
	sync.RWMutex
	streamer  string
	channelId string
	viewers   int
	Users     *odered.OrderedSet[string]
	Emotes    *emote.Store
	Badges    *badge.Store
}

func NewChannel(twitchApi *helix.Client, streamerName string) (*Channel, error) {
	channelId, err := getChannelIdFromId(twitchApi, streamerName)
	if err != nil {
		return nil, err
	}
	badges, err := twitchApi.GetChannelChatBadges(
		&helix.GetChatBadgeParams{
			BroadcasterID: channelId,
		})
	if err != nil {
		return nil, err
	}

	bttvChannelEmotes, err := emote.GetChannelEmotesBTTV(channelId)
	if err != nil {
		return nil, err
	}
	ffzChannelEmotes, err := emote.GetChannelEmotesFFZ(channelId)
	if err != nil {
		return nil, err
	}
	emoteStore := emote.NewStore()
	badgeStore := badge.NewStore()

	emoteStore.SetEmotes(bttvChannelEmotes)
	emoteStore.SetEmotes(ffzChannelEmotes)
	badgeStore.SetBadges(badges.Data.Badges)

	return &Channel{
		RWMutex:   sync.RWMutex{},
		streamer:  streamerName,
		channelId: channelId,
		viewers:   0,
		Users:     odered.NewSet[string](),
		Emotes:    emoteStore,
		Badges:    badgeStore,
	}, nil

}
func (ch *Channel) StreamerId() string {
	return ch.streamer
}

func (ch *Channel) ChannelId() string {
	return ch.channelId
}

func (ch *Channel) UserList() []string {
	return ch.Users.Values()
}

func (ch *Channel) AddUser(userName string) bool {
	if userName == "" {
		return false
	}
	if ch.Users.Exist(userName) {
		return false
	}
	return !ch.Users.Set(userName)
}

func (ch *Channel) RemoveUser(userName string) bool {
	return ch.Users.Delete(userName)
}

func (ch *Channel) ViewerCount() int {
	ch.RLock()
	defer ch.RUnlock()
	return ch.viewers
}

func (ch *Channel) UpdateViewerCount(viewers int) {
	ch.Lock()
	ch.viewers = viewers
	ch.Unlock()
}

func getChannelIdFromId(twitchApi *helix.Client, loginId string) (string, error) {
	users, err := twitchApi.GetUsers(&helix.UsersParams{
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
