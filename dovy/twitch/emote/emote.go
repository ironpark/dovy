package emote

import (
	"fmt"
	"sync"
)

type Emote interface {
	GetCode() string
	URL1X() string
}

type Store struct {
	emotes map[string]string
	sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		emotes:  map[string]string{},
		RWMutex: sync.RWMutex{},
	}
}
func (cm *Store) SetEmotes(emotes []Emote) {
	for _, emote := range emotes {
		cm.Lock()
		cm.emotes[emote.GetCode()] = emote.URL1X()
		cm.Unlock()
	}
}

func (cm *Store) GetEmote(code string) string {
	cm.RLock()
	defer cm.RUnlock()
	return cm.emotes[code]
}

func ImgTag(src string) string {
	return fmt.Sprintf("<img src=\"%s\" class=\"emote\"/>", src)
}
