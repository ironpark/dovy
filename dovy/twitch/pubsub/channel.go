package pubsub

import "fmt"

func TopicBits(channelId string) string {
	return fmt.Sprintf("channel-bits-events-v1.%s", channelId)
}

func TopicStreamInfo(channelId string) string {
	return fmt.Sprintf("video-playback-by-id.%s", channelId)
}
