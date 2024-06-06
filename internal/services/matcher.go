package services

import (
	"container/list"
	"sync"
	"time"

	"github.com/beevik/guid"
)

type Matcher struct {
	queue *MatcherQueue             // очередь игроков, кому нужна комната
	rooms map[string]*GameRoomGroup // комнаты, сгруппированные по играм
}

type MatcherRoom struct {
	players []MatcherPlayer
}

type MatcherPlayer struct {
	playerId guid.Guid
}

type GameRoomGroup struct {
	playersCount int
	gameID       string
	rooms        []MatcherRoom
}

type MatcherQueue struct {
	list *list.List
	lock sync.RWMutex
}

// запрос игрока на комнату
type RoomQuery struct {
	PlayerID guid.Guid
	GameID   string
}

func New() *Matcher {
	m := &Matcher{
		queue: &MatcherQueue{
			list: list.New(),
		},
	}

	// задача на распределение игроков по комнатам
	go func() {
		for {
			m.doMatching()
			time.Sleep(time.Microsecond * 100)
		}
	}()

	return m
}

func (m *Matcher) doMatching() {
	if m.queue.list.Len() == 0 {
		return
	}

	// берем последовательно игроков и добавляем в комнаты, если есть
	// комнаты сгруппированы по типу игры (в мапе)
}

// добавляет запрос на комнату, только если такого еще нет
func (m *Matcher) CheckAndAdd(q RoomQuery) bool {
	m.queue.lock.Lock()
	defer m.queue.lock.Unlock()

	for e := m.queue.list.Front(); e != nil; e = e.Next() {
		if e.Value.(RoomQuery) == q {
			return false
		}
	}

	m.queue.list.PushBack(q)

	return true
}
