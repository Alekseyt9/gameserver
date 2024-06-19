package services

import (
	"container/list"
	"context"
	"gameserver/internal/services/model"
	"gameserver/internal/services/store"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Matcher struct {
	playerManager *PlayerManager
	gameManager   *GameManager
	queue         *MatcherQueue             // очередь игроков, кому нужна комната.
	rooms         map[string]*GameRoomGroup // комнаты, сгруппированные по играм.
	store         store.Store
	log           *slog.Logger
}

type GameRoomGroup struct {
	PlayersCount int
	Rooms        []*model.MatcherRoom
}

type MatcherQueue struct {
	List *list.List
	Lock sync.RWMutex
}

const (
	taskInterval = 500
)

func NewMatcher(store store.Store, pm *PlayerManager, gm *GameManager, log *slog.Logger) (*Matcher, error) {
	m := &Matcher{
		queue: &MatcherQueue{
			List: list.New(),
		},
		store:         store,
		playerManager: pm,
		gameManager:   gm,
		log:           log,
	}

	ctx := context.Background()
	err := m.loadWaitingRooms(ctx)
	if err != nil {
		return nil, err
	}

	// задача на распределение игроков по комнатам.
	go func() {
		ctx = context.Background()
		for {
			err = m.doMatching(ctx)
			if err != nil {
				m.log.Error("m.doMatching error", err)
			}
			time.Sleep(time.Microsecond * taskInterval)
		}
	}()

	return m, nil
}

// загрузка комнат в ожидании игроков из базы.
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
				PlayersCount: m.gameManager.GetGameInfo(r.GameID).PlayerCount,
				Rooms:        make([]*model.MatcherRoom, 0),
			}
			m.rooms[r.GameID] = rg
		}
		rg.Rooms = append(rg.Rooms, r)
	}

	return nil
}

// добавляет запрос на комнату, только если такого еще нет.
func (m *Matcher) CheckAndAdd(q model.RoomQuery) bool {
	m.queue.Lock.Lock()
	defer m.queue.Lock.Unlock()

	for e := m.queue.List.Front(); e != nil; e = e.Next() {
		if e.Value.(model.RoomQuery) == q {
			return false
		}
	}

	m.queue.List.PushBack(q)

	return true
}

func (m *Matcher) doMatching(ctx context.Context) error {
	s := m.queueToSlice()
	if len(s) == 0 {
		return nil
	}

	rooms, err := m.processRooms(s)
	if err != nil {
		return err
	}

	msgs := m.getStartGameMessages()

	// сохраняем все изменения в базу.
	err = m.store.CreateOrUpdateRooms(ctx, rooms)
	if err != nil {
		return err
	}

	// удаляем заполненные комнаты, сбрасываем состояние игроков.
	for g, v := range m.rooms {
		rs := make([]*model.MatcherRoom, 0)
		for _, r := range v.Rooms { // создаем список комнат в ожидании и оставляем только их.
			if len(r.Players) < v.PlayersCount {
				r.IsNew = false
				r.StatusChanged = false
				rs = append(rs, r)
				for _, p := range r.Players { // так как уже в базе - то не новые.
					p.IsNew = false
				}
			}
		}
		m.rooms[g].Rooms = rs
	}

	// рассылаем сообщения о старте игры игрокам (всем игрокам комнат, котрые перешли в режим игры).
	m.sendMessages(msgs)

	return nil
}

func (m *Matcher) processRooms(s []model.RoomQuery) ([]*model.MatcherRoom, error) {
	for _, l := range s {
		rg, ok := m.rooms[l.GameID]
		if !ok {
			rg = &GameRoomGroup{
				Rooms:        make([]*model.MatcherRoom, 0),
				PlayersCount: m.gameManager.GetGameInfo(l.GameID).PlayerCount,
			}
			m.rooms[l.GameID] = rg
		}

		// получаем первую комнату, которая не заполнена.
		var wr *model.MatcherRoom
		for _, r := range rg.Rooms {
			if len(r.Players) < rg.PlayersCount {
				wr = r
				break
			}
		}

		if wr == nil { // создаем новую комнату.
			wr = &model.MatcherRoom{
				ID:      uuid.New(),
				Players: make([]*model.MatcherPlayer, 0),
				IsNew:   true,
				Status:  "wait",
				GameID:  l.GameID,
			}
			rg.Rooms = append(rg.Rooms, wr)
		}

		wr.Players = append(wr.Players, &model.MatcherPlayer{
			PlayerID: l.PlayerID,
			IsNew:    true,
		})
		if len(wr.Players) == rg.PlayersCount {
			wr.Status = "game"
			wr.StatusChanged = true
			err := m.initGame(wr)
			if err != nil {
				return nil, err
			}
		}
	}

	rooms := make([]*model.MatcherRoom, 0)

	for _, v := range m.rooms {
		rooms = append(rooms, v.Rooms...)
	}

	return rooms, nil
}

// инициализируем игру.
func (m *Matcher) initGame(r *model.MatcherRoom) error {
	state, err := m.gameManager.Init(r.GameID, r.Players)
	if err != nil {
		return err
	}
	r.State = state

	return nil
}

// рассылаем сообщения игрокам про старт игры (кто онлайн).
func (m *Matcher) sendMessages(msgs []model.SendMessage) {
	for _, msg := range msgs {
		chp := m.playerManager.GetChan(msg.PlayerID)
		if chp != nil {
			*chp <- msg
		}
	}
}

// получаем сообщения для рассылки игрокам.
func (m *Matcher) getStartGameMessages() []model.SendMessage {
	res := make([]model.SendMessage, 0)

	// только комнаты, которые изменили состояние и перешли в игру.
	for _, g := range m.rooms {
		for _, r := range g.Rooms {
			if r.StatusChanged {
				for _, p := range r.Players {
					res = append(res, model.NewSendMessage(
						p.PlayerID,
						&r.ID,
						model.CreateStartGameMsg(m.gameManager.GetGameInfo(r.GameID).ContentURL),
					))
				}
			}
		}
	}

	return res
}

// перемещаем из очереди все элементы в слайс, чтобы не блокировать очередь надолго, очередь очищается.
func (m *Matcher) queueToSlice() []model.RoomQuery {
	m.queue.Lock.Lock()
	defer m.queue.Lock.Unlock()

	s := make([]model.RoomQuery, 0, m.queue.List.Len())
	l := m.queue.List

	for e := l.Front(); e != nil; e = e.Next() {
		s = append(s, e.Value.(model.RoomQuery))
	}

	l.Init()

	return s
}
