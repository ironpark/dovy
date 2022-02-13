package twitch

import odered "dovey/pkg/odered"

type ConnectionManager struct {
	channels *odered.OrderedMap[string, Channel]
}

func New() *ConnectionManager {
	return &ConnectionManager{
		odered.NewMap[string, Channel](),
	}
}

func (cm *ConnectionManager) Connect(channel string) {
	cm.channels.Set(channel)
}
