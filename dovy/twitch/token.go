package twitch

import (
	"fmt"
	"github.com/nicklaw5/helix/v2"
	"math/rand"
	"net/http"
	"strings"
	"sync"
)

const TwitchClientId = "xxt39lhzhqa7z2wwfdzjkdrtq4cj21"

type TokenRecvCallback func(token string)
type TokenReceiver struct {
	server   *http.Server
	mux      *http.ServeMux
	callback TokenRecvCallback
	lock     sync.RWMutex
}

func NewTokenReceiver() *TokenReceiver {
	mux := http.NewServeMux()
	tokenReceiver := &TokenReceiver{
		server: &http.Server{Addr: fmt.Sprintf("localhost:%d", 53324), Handler: mux},
		mux:    mux,
	}
	mux.HandleFunc("/authorize", tokenReceiver.authorize)
	mux.HandleFunc("/token", tokenReceiver.token)
	return tokenReceiver
}

func (tr *TokenReceiver) SetTokenRecvCallback(callback TokenRecvCallback) {
	tr.lock.Lock()
	tr.callback = callback
	tr.lock.Unlock()
}

func (tr *TokenReceiver) authorize(w http.ResponseWriter, req *http.Request) {
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

func (tr *TokenReceiver) token(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	accessToken := query.Get("access_token")
	tr.lock.RLock()
	defer tr.lock.RUnlock()
	if tr.callback == nil {
		return
	}
	tr.callback(accessToken)
}

func GetAuthorizationURL(scope []string, forceVerify bool) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyz")
	randomRunes := make([]rune, 15)
	for i := range randomRunes {
		randomRunes[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	state := string(randomRunes)
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
