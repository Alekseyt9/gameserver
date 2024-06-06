package services

import (
	"gameserver/internal/services/model"
	"gameserver/internal/services/store"
	"sync"

	"github.com/beevik/guid"
)

type RoomManager struct {
	store       store.Store
	interactors *InteractorMap
	gameManager *GameManager
}

type InteractorMap struct {
	data map[guid.Guid]*RoomInteractor
	lock sync.RWMutex
}

type RoomInfo struct {
}

type RoomInteractor struct {
	RecieveChan chan RecieveRoomMsg // канал для обработки входящих через websocket сообщений одним воркером комнаты
	SendChan    chan SendRooomMsg   // канал для рассылки сообщений игрокам комнаты через werbsocket
}

type RecieveRoomMsg struct {
}

type SendRooomMsg struct {
}

func (r *RoomManager) GetRoomInteractor(roomID guid.Guid) *RoomInteractor {
	r.interactors.lock.Lock()
	defer r.interactors.lock.Unlock()
	x, ok := r.interactors.data[roomID]
	if !ok {

		x = &RoomInteractor{
			RecieveChan: make(chan RecieveRoomMsg),
			SendChan:    make(chan SendRooomMsg),
		}

		// создаем воркер для комнаты
		go func() {
			for m := range x.RecieveChan {
				r.gameManager.Process(m)
			}
		}()

		r.interactors.data[roomID] = x
	}
	return x
}

// подключение к существующей комнате или создание комнаты
// подключаться нужно каждый раз при коннекте игрока
func (r *RoomManager) PlayerConnect(msg *model.ClientMsg) {
	room := r.getExistingRoom(msg.PlayerID, msg.GameType)
	if room != nil {
		// комната уже есть,
		// если комната в режиме игры - стартуем игру
		// !!! TODO если в режиме ожидания
	} else {
		// !!! TODO ставим в очередь в Matcher
	}
}

// выход из комнаты (выйти можно только один раз)
func (r *RoomManager) PlayerQuit(msg *model.ClientMsg) {

}

func (r *RoomManager) getExistingRoom(userId guid.Guid, gameType string) *RoomInfo {
	return nil
}

/*
	!TODO

-- Получает существующую комнату для игрока или создает новую
-- Существующую: если игра с таким типом для игрока уже есть или игры нет, но игрок подключается к свободной комнате
-- Новую: если нет свободной комнаты для игрока и нет существующей игры игры
*/
func (r *RoomManager) getOrCreateRoom(userId guid.Guid, gameType string) {
	// запрос в хранилище по игроку и типу игры
	// ? нужна ли комната в базе, котоорая ожидает игроков?
	// наверное да, тк чтобы игроки могли подключаться, когда другие офлайн
}
