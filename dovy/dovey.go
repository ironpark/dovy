package dovy

import (
	"context"
	"dovey/dovy/twitch"
	"dovey/dovy/twitch/chat"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os/exec"
	gort "runtime"
	"sync"
)

const TwitchClientId = "xxt39lhzhqa7z2wwfdzjkdrtq4cj21"

type Dovy struct {
	appCtx      context.Context
	accessToken string
	scope       []string
	cm          *chat.ChatManager
	lock        sync.Mutex
}

func NewDovey(ctx context.Context) *Dovy {
	cm := chat.NewChatManager(TwitchClientId)
	dovy := &Dovy{
		appCtx: ctx,
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
		cm: cm,
	}
	cm.OnMsg(func(message chat.Message) {
		runtime.EventsEmit(dovy.appCtx, "chat.stream", message)
	})
	return dovy
}

func (dov Dovy) OpenAuthorization() {
	authUrl := twitch.GetAuthorizationURL(dov.scope, true)
	fmt.Println(authUrl)
	switch gort.GOOS {
	case "darwin":
		exec.Command("open", authUrl).Run()
	case "windows":
		exec.Command("start", "/max", authUrl).Run()
	}
}

func (dov *Dovy) SendChatMessage(msg string) {
	fmt.Println("SendChatMessage", "woowakgood", msg)
}
