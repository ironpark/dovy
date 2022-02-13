package twitch

import odered "dovey/pkg/odered"

type Channel struct {
	streamer  string
	channelId string
	viewers   int
	users     odered.OrderedSet[string]
}
