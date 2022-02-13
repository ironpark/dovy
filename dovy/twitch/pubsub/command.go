package pubsub

import (
	"fmt"
	"strings"
)

func listen(token, nonce string, topics ...string) []byte {
	for i, topic := range topics {
		topics[i] = fmt.Sprintf(`"%s"`, topic)
	}
	topicsStr := strings.Join(topics, ",")
	cmd := fmt.Sprintf(`{"type":"LISTEN","nonce":"%s","data":{"topics":[%s],"auth_token":"%s"}}`, nonce, topicsStr, token)
	return []byte(cmd)
}
