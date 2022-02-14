package twitch

import (
	"dovey/dovy/twitch/badge"
	"dovey/dovy/twitch/emote"
	odered "dovey/pkg/odered"
	"sync"
)

type Channel struct {
	sync.RWMutex
	streamer   string
	channelId  string
	viewers    int
	users      *odered.OrderedSet[string]
	emoteStore *emote.Store
	badgeStore *badge.Store
}
