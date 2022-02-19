package pubsub

import (
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

const endpoint = "wss://pubsub-edge.twitch.tv"

type PubSub struct {
	conn *websocket.Conn
}

func New() (*PubSub, error) {
	conn, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		return nil, err
	}
	pubsub := &PubSub{conn: conn}
	go pubsub.pingLoop()
	return pubsub, nil
}

func (ps *PubSub) pingLoop() error {
	t := time.NewTicker(time.Minute)
	defer t.Stop()
	for range t.C {
		err := ps.conn.WriteMessage(websocket.TextMessage, []byte(`{"type": "PING"}`))
		if err != nil {
			return err
		}
	}
	return nil
}
func (ps *PubSub) ReadLoop() error {
	for {
		t, p, err := ps.conn.ReadMessage()
		if err != nil {
			return err
		}
		switch t {
		case websocket.TextMessage:
			fmt.Println(string(p))
		case websocket.BinaryMessage:
		case websocket.PingMessage:
		case websocket.PongMessage:
		}

	}
}

func (ps *PubSub) Listen(token, nonce string, topics ...string) {
	c := listen(token, nonce, topics...)
	ps.conn.WriteMessage(websocket.TextMessage, c)
}
