package dovy

import (
	"context"
	"dovey/dovy/twitch/chat"
	"fmt"
	"github.com/nicklaw5/helix/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"math/rand"
	"net/http"
	"os/exec"
	gort "runtime"
	"strings"
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

func (dov *Dovy) token(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	accessToken := query.Get("access_token")
	//scope := query.Get("scope")
	//state := query.Get("state")

	dov.lock.Lock()
	dov.accessToken = accessToken
	dov.lock.Unlock()

	err := dov.cm.Initialize(dov.accessToken)
	if err != nil {
		fmt.Println(err)
		return
	}

	dov.cm.Connect("woowakgood")
}

func (dov *Dovy) authorize(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(`
<!doctype html>
<html>
<head>
	<title>Dovy</title>
</head>
<body>
	<div id="result">
		<p>도비에게 인증 토큰을 전달하는 중입니다 잠시만 기다려주세요.</p>
		<p>Passing auth token to dovy, please wait.</p>
	</div>
	<script>
		var resultContainer = document.getElementById("result")
		var xmlHttp = new XMLHttpRequest()
		xmlHttp.onreadystatechange = function () {
		  if (this.status == 200 && this.readyState == this.DONE) {
			resultContainer.innerHTML =  "<p>도비는 이제 일할 준비가 완료되었습니다. 이 창은 이제 꺼주셔도 됩니다.</p>"
			resultContainer.innerHTML += "<p>Dovy is now ready to go to work. You can now close this window.</p>"
			window.close();
		  }
		}
		xmlHttp.open('GET', '/token?' + document.location.hash.substring(1), true)
		xmlHttp.send();
	</script>
</body>
</html>`))
}

func (dov *Dovy) Serve() error {
	http.HandleFunc("/authorize", dov.authorize)
	http.HandleFunc("/token", dov.token)
	return http.ListenAndServe(":53324", nil)
}

func (dov Dovy) OpenAuthorization() {
	authUrl := dov.GetAuthorizationURL()
	fmt.Println(authUrl)
	switch gort.GOOS {
	case "darwin":
		exec.Command("open", authUrl).Run()
	case "windows":
		exec.Command("start", "/max", authUrl).Run()
	}
}

func (dov Dovy) GetAuthorizationURL() string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, 15)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return getAuthorizationURL(dov.scope, string(b), true)
}

func (dov *Dovy) SendChatMessage(msg string) {
	fmt.Println("SendChatMessage", "woowakgood", msg)
}

func getAuthorizationURL(scope []string, state string, forceVerify bool) string {
	url := helix.AuthBaseURL + "/authorize"
	url += "?response_type=" + "token"
	url += "&client_id=" + TwitchClientId
	url += "&redirect_uri=" + "http://localhost:53324/authorize"
	url += "&state=" + state
	if forceVerify {
		url += "&force_verify=true"
	}
	if len(scope) != 0 {
		url += "&scope=" + strings.Join(scope, "%20")
	}
	return url
}
