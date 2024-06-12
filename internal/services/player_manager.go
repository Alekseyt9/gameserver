package services

import (
	"gameserver/internal/services/model"
	"gameserver/internal/services/store"
	"sync"

	"github.com/beevik/guid"
)

type PlayerManager struct {
	store    store.Store
	chanMap  map[guid.Guid]chan model.SendMessage
	chanLock sync.RWMutex
}

func NewPlayerManager(store store.Store) *PlayerManager {
	return &PlayerManager{
		store:   store,
		chanMap: make(map[guid.Guid]chan model.SendMessage),
	}
}

// получить канал для свази с игроком
func (m *PlayerManager) GetOrCreateChan(palyerID guid.Guid) chan model.SendMessage {
	m.chanLock.Lock()
	defer m.chanLock.Unlock()

	ch, ok := m.chanMap[palyerID]
	if !ok {
		ch = make(chan model.SendMessage, 100)
		m.chanMap[palyerID] = ch
	}
	return ch
}

// получить канал для свази с игроком
func (m *PlayerManager) GetChan(palyerID guid.Guid) *chan model.SendMessage {
	m.chanLock.RLock()
	defer m.chanLock.RUnlock()

	ch, ok := m.chanMap[palyerID]
	if !ok {
		return nil
	}
	return &ch
}

// удалить канал для свзяи с игроком
func (m *PlayerManager) DeleteChan(palyerID guid.Guid) {
	m.chanLock.Lock()
	defer m.chanLock.Unlock()

	ch, ok := m.chanMap[palyerID]
	if ok {
		close(ch)
		delete(m.chanMap, palyerID)
	}
}
