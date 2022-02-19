package badge

import (
	"fmt"
	"github.com/nicklaw5/helix/v2"
	"sync"
)

type Store struct {
	badges map[string]string
	sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		badges:  map[string]string{},
		RWMutex: sync.RWMutex{},
	}
}
func (cm *Store) SetBadges(badges []helix.ChatBadge) {
	for _, badge := range badges {
		for _, version := range badge.Versions {
			cm.Lock()
			cm.badges[fmt.Sprintf("%s#%s", badge.SetID, version.ID)] = version.ImageUrl1x
			cm.Unlock()
		}
	}
}

func (cm *Store) SetBadge(badgeId string, versionId int, imgUri string) {
	cm.Lock()
	defer cm.Unlock()
	cm.badges[fmt.Sprintf("%s#%d", badgeId, versionId)] = imgUri
}

func (cm *Store) GetBadgeImage(badgeId string, versionId int) string {
	cm.RLock()
	defer cm.RUnlock()
	return cm.badges[fmt.Sprintf("%s#%d", badgeId, versionId)]
}
