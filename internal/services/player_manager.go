package services

import (
	"gameserver/internal/services/model"
	"gameserver/internal/services/store"

	"github.com/beevik/guid"
)

type PlayerManager struct {
	store   store.Store
	chanMap map[guid.Guid]chan model.GameMsg
}

func NewPlayerManager(store store.Store) *PlayerManager {
	return &PlayerManager{
		store:   store,
		chanMap: make(map[guid.Guid]chan model.GameMsg),
	}
}

/* в хендлере
// создать нового пользователя (анонимного)
func (m *PlayerManager) CreateUser() *model.Player {
	ctx := context.Background()
	p := model.Player{
		ID: *guid.New(),
		Name: ,
	}
	m.store.CreateUser(ctx, )
}
*/

// получить канал для свази с игроком
func (m *PlayerManager) GetOrCreateChan(palyerID guid.Guid) chan model.GameMsg {
	ch, ok := m.chanMap[palyerID]
	if !ok {
		ch = make(chan model.GameMsg, 100)
		m.chanMap[palyerID] = ch
	}
	return ch
}

// удалить канал для свзяи с игроком
func (m *PlayerManager) DeleteChan(palyerID guid.Guid) {
	ch, ok := m.chanMap[palyerID]
	if ok {
		close(ch)
		delete(m.chanMap, palyerID)
	}
}
