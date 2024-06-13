package store

import (
	"context"
	"database/sql"
	"errors"
	"gameserver/internal/services/model"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // needs for init
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
)

type DBStore struct {
	conn *sql.DB
}

func NewDBStore(connString string) (Store, error) {
	conn, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}

	err = bootstrap(connString)
	if err != nil {
		return nil, err
	}

	return &DBStore{
		conn: conn,
	}, nil
}

func (s *DBStore) CreatePlayer(ctx context.Context, p *model.Player) error {
	_, err := s.conn.ExecContext(ctx, `
		insert into Players(Id, Name)
		values ($1, $2)
	`, p.ID, p.Name)
	return err
}

func (s *DBStore) GetPlayer(ctx context.Context, playerID uuid.UUID) (*model.Player, error) {
	row := s.conn.QueryRowContext(ctx, `SELECT Name FROM Players WHERE Id = $1`, playerID)
	var name string
	err := row.Scan(&name)

	if err != nil {
		return nil, err
	}

	return &model.Player{
		ID:   playerID,
		Name: name,
	}, nil
}

func (s *DBStore) GetRoom(ctx context.Context, gameID string, playerID uuid.UUID) (*model.Room, error) {
	row := s.conn.QueryRowContext(ctx, `
		SELECT r.Id, r.State
		FROM Rooms r 
		join RoomPlayers rp on rp.RoomID = r.Id 
		WHERE r.GameId = $1 and rp.PlayerId = $2 and not rp.IsQuit
		`, gameID, playerID)

	var id uuid.UUID
	var state string
	err := row.Scan(&id, &state)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &model.Room{ID: id, State: state}, nil
}

func (s *DBStore) SetRoomState(ctx context.Context, roomID uuid.UUID, state string) error {
	_, err := s.conn.ExecContext(ctx, `
		update Rooms
		set State = $1
		where Id = $2
	`, roomID, state)
	return err
}

func (s *DBStore) CreateOrUpdateRooms(ctx context.Context, rooms []*model.MatcherRoom) error {
	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck //defer

	stmtRoomInsert, err := tx.PrepareContext(ctx, `
		insert into Rooms(Id, GameId, Status, State)
		values ($1, $2, $3, $4)
	`)
	if err != nil {
		return err
	}
	defer stmtRoomInsert.Close()

	stmtRoomUpdate, err := tx.PrepareContext(ctx, `
		update Rooms
		set Status = $1, State = $3
		where Id = $2
	`)
	if err != nil {
		return err
	}
	defer stmtRoomUpdate.Close()

	stmtRoomPlayerInsert, err := tx.PrepareContext(ctx, `
		insert into RoomPlayers(PlayerId, RoomId)
		values ($1, $2)
	`)
	if err != nil {
		return err
	}
	defer stmtRoomPlayerInsert.Close()

	for _, r := range rooms {
		if r.IsNew {
			_, err = stmtRoomInsert.ExecContext(ctx, uuid.New(), r.GameID, r.Status, r.State)
			if err != nil {
				return err
			}
		} else {
			if r.StatusChanged {
				_, err = stmtRoomUpdate.ExecContext(ctx, r.Status, r.State)
				if err != nil {
					return err
				}
			}
		}

		for _, p := range r.Players {
			if p.IsNew {
				_, err = stmtRoomPlayerInsert.ExecContext(ctx, p.PlayerID, r.ID)
				if err != nil {
					return err
				}
			}
		}
	}

	return tx.Commit()
}

func (s *DBStore) LoadWaitingRooms(ctx context.Context) ([]*model.MatcherRoom, error) {
	res := make([]*model.MatcherRoom, 0)
	var rows *sql.Rows

	// одним запросом загружаем комнаты и игроков
	rows, err := s.conn.QueryContext(ctx, `
		select r.Id, r.GameId, r.Status, rp.PlayerId
		from Rooms r
		left join RoomPlayers rp on rp.RoomId = r.Id
		where r.Status = 'wait'
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		curRoomId uuid.UUID
		room      *model.MatcherRoom
		roomID    uuid.UUID
		gameID    string
		status    string
		playerID  *uuid.UUID
	)

	for rows.Next() {
		if err = rows.Scan(&roomID, &gameID, &status, &playerID); err != nil {
			return nil, err
		}
		if curRoomId != roomID {
			room = &model.MatcherRoom{
				ID:      roomID,
				Players: make([]*model.MatcherPlayer, 0),
				Status:  status,
				GameID:  gameID,
			}
			curRoomId = room.ID
			res = append(res, room)
		} else {
			if playerID != nil {
				room.Players = append(room.Players, &model.MatcherPlayer{
					PlayerID: *playerID,
				})
			}
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *DBStore) MarkDropRoomPlayer(ctx context.Context, roomID uuid.UUID, playerID uuid.UUID) error {
	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck //defer

	_, err = tx.ExecContext(ctx, `
		update RoomPlayers
		set IsQuit = true
		where RoomId = $1 and PlayerId = $2
	`, roomID, playerID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func bootstrap(connString string) error {
	mPath := getMigrationPath()
	m, err := migrate.New(mPath, connString)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func getMigrationPath() string {
	_, currentFilePath, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFilePath)
	migrationsPath := filepath.Join(currentDir, "migrations")
	migrationsPath = filepath.ToSlash(migrationsPath)
	return "file://" + migrationsPath
}
