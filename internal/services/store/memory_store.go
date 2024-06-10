package store

import (
	"context"
	"gameserver/internal/services/model"

	"github.com/beevik/guid"
)

// для тестов
type MemStore struct {
	rooms       map[guid.Guid]*MSRoom
	players     map[guid.Guid]*MSPlayer
	roomPlayers []MSRoomPlayer
}

type MSRoom struct {
	ID     guid.Guid
	GameID string
	State  string
	Status string
}

type MSPlayer struct {
	ID   guid.Guid
	Name string
}

type MSRoomPlayer struct {
	PlayerID guid.Guid
	RoomID   guid.Guid
}

func NewMemStore() *MemStore {
	return &MemStore{
		rooms:       make(map[guid.Guid]*MSRoom, 0),
		players:     make(map[guid.Guid]*MSPlayer, 0),
		roomPlayers: make([]MSRoomPlayer, 0),
	}
}

func (s *MemStore) GetPlayer(ctx context.Context, id guid.Guid) (*model.Player, error) {
	v, ok := s.players[id]
	if ok {
		return &model.Player{
			ID:   v.ID,
			Name: v.Name,
		}, nil
	}
	return nil, nil
}

func (s *MemStore) CreatePlayer(ctx context.Context, p *model.Player) error {
	s.players[p.ID] = &MSPlayer{
		ID:   p.ID,
		Name: p.Name,
	}
	return nil
}

func (s *MemStore) GetRoom(ctx context.Context, playerID guid.Guid, gameID string) (*model.Room, error) {
	var res *model.Room

exit:
	for _, p := range s.roomPlayers {
		if p.PlayerID == playerID {
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

func (s *MemStore) SetRoomState(ctx context.Context, id guid.Guid, state string) error {
	r, ok := s.rooms[id]
	if ok {
		r.State = state
	}
	return nil
}

func (s *MemStore) CreateOrUpdateRooms(ctx context.Context, rooms []model.MatcherRoom) error {
	for _, r := range rooms {
		if r.IsNew {
			s.rooms[r.ID] = &MSRoom{
				ID:     r.ID,
				GameID: r.GameID,
				State:  "",
				Status: r.Status,
			}
		} else {
			old, ok := s.rooms[r.ID]
			if ok {
				old.Status = r.Status
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

func (s *MemStore) LoadWaitingRooms(ctx context.Context) ([]model.MatcherRoom, error) {
	res := make([]model.MatcherRoom, 0)

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
			res = append(res, *r)
		}
	}

	return res, nil
}
