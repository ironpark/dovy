package emote

import (
	"fmt"
	"sync"
)

type Emote interface {
	GetCode() string
	URL1X() string
	GetImgTag() string
}

type Store struct {
	emotes map[string]Emote
	sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		emotes:  map[string]Emote{},
		RWMutex: sync.RWMutex{},
	}
}
func (cm *Store) SetEmotes(emotes []Emote) {
	if emotes == nil {
		return
	}
	for _, emote := range emotes {
		cm.Lock()
		cm.emotes[emote.GetCode()] = emote
		cm.Unlock()
	}
}

func (cm *Store) GetEmote(code string) Emote {
	cm.RLock()
	defer cm.RUnlock()
	return cm.emotes[code]
}

func ImgTag(src, code string) string {
	return fmt.Sprintf("<img src=\"%s\" class=\"emote\" alt=\"%s\"/>", src, code)
}
