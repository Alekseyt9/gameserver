package services

import (
	"container/list"
	"gameserver/internal/services/model"
	"gameserver/internal/services/store"
	"sync"
	"time"

	"github.com/beevik/guid"
)

type Matcher struct {
	queue *MatcherQueue             // очередь игроков, кому нужна комната
	rooms map[string]*GameRoomGroup // комнаты, сгруппированные по играм
	store store.Store
}

type GameRoomGroup struct {
	playersCount int
	rooms        []model.MatcherRoom
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
			time.Sleep(time.Microsecond * 500)
		}
	}()

	return m
}

func (m *Matcher) doMatching() error {
	s := m.queueToSlice()

	for _, l := range s {
		rg := m.rooms[l.GameID]

		// получаем первую комнату, которая не заполнена
		var wr *model.MatcherRoom
		for _, r := range rg.rooms {
			if len(r.Players) < rg.playersCount {
				wr = &r
				break
			}
		}

		if wr == nil { // создаем новую комнату
			wr = &model.MatcherRoom{
				Players: make([]model.MatcherPlayer, rg.playersCount),
				IsNew:   true,
				Status:  "wait",
			}
			rg.rooms = append(rg.rooms, *wr)
		}

		wr.Players = append(wr.Players, model.MatcherPlayer{
			PlayerId: l.PlayerID,
			IsNew:    true,
		})
		if len(wr.Players) == rg.playersCount {
			wr.Status = "game"
		}
	}

	// сохраняем все изменения
	rooms := make([]model.MatcherRoom, 0)

	for _, v := range m.rooms {
		for _, r := range v.rooms {
			rooms = append(rooms, r)
		}
	}

	err := m.store.CreateOrUpdateRooms(rooms)
	if err != nil {
		return nil
	}

	// удаляем заполненные комнаты, сбрасываем состояние игроков
	for g, v := range m.rooms {
		rs := make([]model.MatcherRoom, 0)
		for _, r := range v.rooms { // создаем список комнат в ожидании и оставляем только их
			if len(r.Players) < v.playersCount {
				r.IsNew = false
				rs = append(rs, r)
				for _, p := range r.Players { // так как уже в базе - то не новые
					p.IsNew = false
				}
			}
		}
		m.rooms[g].rooms = rs
	}

	return nil
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

// перемещаем из очереди все элементы в слайс, чтобы не блокировать очередь надолго, очередь очищается
func (m *Matcher) queueToSlice() []RoomQuery {
	m.queue.lock.Lock()
	defer m.queue.lock.Unlock()
	s := make([]RoomQuery, m.queue.list.Len())
	l := m.queue.list

	for e := l.Front(); e != nil; {
		next := e.Next()
		s = append(s, e.Value.(RoomQuery))
		l.Remove(e)
		e = next
	}

	return s
}
