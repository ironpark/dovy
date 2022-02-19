package dovy

import (
	"context"
	"dovey/dovy/twitch"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os/exec"
	gort "runtime"
	"sync"
)

const TwitchClientId = "xxt39lhzhqa7z2wwfdzjkdrtq4cj21"

type Dovy struct {
	appCtx        context.Context
	accessToken   string
	scope         []string
	lock          sync.Mutex
	tokenReceiver *twitch.TokenReceiver
	cm            *twitch.ConnectionManager
}

func NewDovey() (*Dovy, error) {
	cm, err := twitch.NewConnectionManager()
	if err != nil {
		return nil, err
	}

	dovy := &Dovy{
		appCtx:      nil,
		accessToken: "",
		scope: []string{
			"channel:moderate",
			"chat:edit",
			"chat:read",
			"whispers:read",
			"whispers:edit",
			"moderator:manage:chat_settings",
			"channel:edit:commercial",
			"channel:manage:broadcast",
			"user:manage:blocked_users",
			"user:read:blocked_users",
			"moderator:manage:banned_users",
			"moderator:read:blocked_terms",
		},
		lock:          sync.Mutex{},
		tokenReceiver: twitch.NewTokenReceiver(),
		cm:            cm,
	}

	dovy.tokenReceiver.SetTokenRecvCallback(func(token string) {
		cm.Initialize(token)
	})

	go dovy.tokenReceiver.Serve()
	return dovy, nil
}

func (dov *Dovy) SetAppContext(ctx context.Context) {
	dov.lock.Lock()
	defer dov.lock.Unlock()
	if dov.appCtx != nil {
		return
	}
	dov.appCtx = ctx
	dov.cm.OnMessage(func(message twitch.Message) {
		runtime.EventsEmit(ctx, "stream.chat", message)
	})
	dov.cm.OnUserEvent(func(channel *twitch.Channel, event string, user string) {
		runtime.EventsEmit(ctx, "stream.user-event", map[string]interface{}{
			"channel": channel.StreamerId(),
			"event":   event,
			"user":    user,
		})
	})
	return
}

func (dov *Dovy) IsAuthorized() bool {
	return dov.cm.IsInitialized()
}

func (dov Dovy) OpenAuthorization() {
	authUrl := twitch.GetAuthorizationURL(dov.scope, true)
	switch gort.GOOS {
	case "darwin":
		exec.Command("open", authUrl).Run()
	case "windows":
		exec.Command("start", "/max", authUrl).Run()
	}
}

func (dov *Dovy) UserList(channelName string) []string {
	channel, err := dov.cm.ConnectOrGet(channelName)
	if err != nil {
		return []string{}
	}
	return channel.UserList()
}

func (dov *Dovy) Connect(channelName string) {
	dov.cm.ConnectOrGet(channelName)
}

func (dov *Dovy) SendChatMessage(channel, msg string) {
	dov.cm.Send(channel, msg)
}
