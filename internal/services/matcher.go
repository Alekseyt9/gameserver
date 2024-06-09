package services

import (
	"container/list"
	"context"
	"gameserver/internal/services/model"
	"gameserver/internal/services/store"
	"sync"
	"time"
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

func New() (*Matcher, error) {
	m := &Matcher{
		queue: &MatcherQueue{
			list: list.New(),
		},
	}

	ctx := context.Background()
	err := m.loadWaitingRooms(ctx)
	if err != nil {
		return nil, err
	}

	// задача на распределение игроков по комнатам
	go func() {
		ctx := context.Background()
		for {
			m.doMatching(ctx)
			time.Sleep(time.Microsecond * 500)
		}
	}()

	return m, nil
}

// загрузка комнат в ожидании игроков из базы
func (m *Matcher) loadWaitingRooms(ctx context.Context) error {
	rooms, err := m.store.LoadWaitingRooms(ctx)
	if err != nil {
		return err
	}

	m.rooms = make(map[string]*GameRoomGroup)
	for _, r := range rooms {
		rg, ok := m.rooms[r.GameID]
		if !ok {
			rg = &GameRoomGroup{
				playersCount: 2, // TODO !! брать из обработчика игры
				rooms:        make([]model.MatcherRoom, 0),
			}
			m.rooms[r.GameID] = rg
		}
		rg.rooms = append(rg.rooms, r)
	}

	return nil
}

// добавляет запрос на комнату, только если такого еще нет
func (m *Matcher) CheckAndAdd(q model.RoomQuery) bool {
	m.queue.lock.Lock()
	defer m.queue.lock.Unlock()

	for e := m.queue.list.Front(); e != nil; e = e.Next() {
		if e.Value.(model.RoomQuery) == q {
			return false
		}
	}

	m.queue.list.PushBack(q)

	return true
}

func (m *Matcher) doMatching(ctx context.Context) error {
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
		rooms = append(rooms, v.rooms...)
	}

	err := m.store.CreateOrUpdateRooms(ctx, rooms)
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

// перемещаем из очереди все элементы в слайс, чтобы не блокировать очередь надолго, очередь очищается
func (m *Matcher) queueToSlice() []model.RoomQuery {
	m.queue.lock.Lock()
	defer m.queue.lock.Unlock()
	s := make([]model.RoomQuery, m.queue.list.Len())
	l := m.queue.list

	for e := l.Front(); e != nil; {
		next := e.Next()
		s = append(s, e.Value.(model.RoomQuery))
		l.Remove(e)
		e = next
	}

	return s
}
