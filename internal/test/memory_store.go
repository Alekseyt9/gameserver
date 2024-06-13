package test

import (
	"context"
	"gameserver/internal/services/model"
	"sync"

	"github.com/google/uuid"
)

// для тестов
type MemStore struct {
	rooms       map[uuid.UUID]*MSRoom
	players     map[uuid.UUID]*MSPlayer
	roomPlayers []MSRoomPlayer
	lock        sync.RWMutex
}

type MSRoom struct {
	ID     uuid.UUID
	GameID string
	State  string
	Status string
}

type MSPlayer struct {
	ID   uuid.UUID
	Name string
}

type MSRoomPlayer struct {
	PlayerID uuid.UUID
	RoomID   uuid.UUID
	IsQuit   bool
}

func NewMemStore() *MemStore {
	return &MemStore{
		rooms:       make(map[uuid.UUID]*MSRoom, 0),
		players:     make(map[uuid.UUID]*MSPlayer, 0),
		roomPlayers: make([]MSRoomPlayer, 0),
	}
}

func (s *MemStore) GetPlayer(ctx context.Context, playerID uuid.UUID) (*model.Player, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	p, ok := s.players[playerID]
	if ok {
		return &model.Player{
			ID:   p.ID,
			Name: p.Name,
		}, nil
	}
	return nil, nil
}

func (s *MemStore) CreatePlayer(ctx context.Context, p *model.Player) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.players[p.ID] = &MSPlayer{
		ID:   p.ID,
		Name: p.Name,
	}
	return nil
}

func (s *MemStore) GetRoom(ctx context.Context, gameID string, playerID uuid.UUID) (*model.Room, error) {
	var res *model.Room
	s.lock.RLock()
	defer s.lock.RUnlock()

exit:
	for _, p := range s.roomPlayers {
		if p.PlayerID == playerID && !p.IsQuit {
			if s.rooms[p.RoomID].GameID == gameID {
				r := s.rooms[p.RoomID]
				res = &model.Room{
					ID:     r.ID,
					State:  r.State,
					Status: r.Status,
				}
				break exit
			}
		}
	}

	return res, nil
}

func (s *MemStore) SetRoomState(ctx context.Context, id uuid.UUID, state string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	r, ok := s.rooms[id]
	if ok {
		r.State = state
	}
	return nil
}

func (s *MemStore) CreateOrUpdateRooms(ctx context.Context, rooms []*model.MatcherRoom) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, r := range rooms {
		if r.IsNew {
			s.rooms[r.ID] = &MSRoom{
				ID:     r.ID,
				GameID: r.GameID,
				State:  r.State,
				Status: r.Status,
			}
		} else {
			if r.StatusChanged {
				old, ok := s.rooms[r.ID]
				if ok {
					old.Status = r.Status
					old.State = r.State
				}
			}
		}

		for _, p := range r.Players {
			if p.IsNew {
				s.roomPlayers = append(s.roomPlayers, MSRoomPlayer{
					PlayerID: p.PlayerID,
					RoomID:   r.ID,
				})
			}
		}
	}

	return nil
}

func (s *MemStore) LoadWaitingRooms(ctx context.Context) ([]*model.MatcherRoom, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	res := make([]*model.MatcherRoom, 0)

	for _, v := range s.rooms {
		if v.Status == "wait" {
			ps := make([]model.MatcherPlayer, 0)
			for _, p := range s.roomPlayers {
				if p.RoomID == v.ID {
					ps = append(ps, model.MatcherPlayer{
						PlayerID: p.PlayerID,
					})
				}
			}

			r := &model.MatcherRoom{
				ID:      v.ID,
				Status:  v.Status,
				Players: ps,
				GameID:  v.GameID,
			}
			res = append(res, r)
		}
	}

	return res, nil
}

func (s *MemStore) MarkDropRoomPlayer(ctx context.Context, roomID uuid.UUID, playerID uuid.UUID) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, r := range s.roomPlayers {
		if r.RoomID == roomID && r.PlayerID == playerID {
			r.IsQuit = true
		}
	}

	/*
		_, ok := s.rooms[roomID]
		if ok {
			prs := make([]MSRoomPlayer, 0)
			for _, p := range s.roomPlayers {
				if !(p.PlayerID == playerID && p.RoomID == roomID) {
					prs = append(prs, p)
				}
			}
			s.roomPlayers = prs

			// удаляем комнату, если она пустая
			hasPlayers := false
			for _, p := range s.roomPlayers {
				if p.RoomID == roomID {
					hasPlayers = true
				}
			}
			if !hasPlayers {
				delete(s.rooms, roomID)
			}
		}*/

	return nil
}
