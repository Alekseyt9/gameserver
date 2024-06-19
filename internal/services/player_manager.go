package services

import (
	"gameserver/internal/services/model"
	"gameserver/internal/services/store"
	"sync"

	"github.com/google/uuid"
)

type PlayerManager struct {
	store    store.Store
	chanMap  map[uuid.UUID]chan model.SendMessage
	chanLock sync.RWMutex
}

func NewPlayerManager(store store.Store) *PlayerManager {
	return &PlayerManager{
		store:   store,
		chanMap: make(map[uuid.UUID]chan model.SendMessage),
	}
}

func (m *PlayerManager) SendToPlayer(playerID uuid.UUID, msg model.SendMessage) error {
	ch := m.GetOrCreateChan(playerID)
	ch <- msg
	return nil
}

// получить канал для свази с игроком.
func (m *PlayerManager) GetOrCreateChan(palyerID uuid.UUID) chan model.SendMessage {
	m.chanLock.Lock()
	defer m.chanLock.Unlock()

	ch, ok := m.chanMap[palyerID]
	if !ok {
		ch = make(chan model.SendMessage, chanBuffer)
		m.chanMap[palyerID] = ch
	}
	return ch
}

// получить канал для свази с игроком.
func (m *PlayerManager) GetChan(palyerID uuid.UUID) *chan model.SendMessage {
	m.chanLock.RLock()
	defer m.chanLock.RUnlock()

	ch, ok := m.chanMap[palyerID]
	if !ok {
		return nil
	}
	return &ch
}

// удалить канал для свзяи с игроком.
func (m *PlayerManager) DeleteChan(palyerID uuid.UUID) {
	m.chanLock.Lock()
	defer m.chanLock.Unlock()

	ch, ok := m.chanMap[palyerID]
	if ok {
		close(ch)
		delete(m.chanMap, palyerID)
	}
}
