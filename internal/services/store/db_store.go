package store

import (
	"context"
	"database/sql"
	"errors"
	"gameserver/internal/services/model"
	"log/slog"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // needs for init
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
)

type DBStore struct {
	conn *sql.DB
	log  *slog.Logger
}

func NewDBStore(connString string, log *slog.Logger) (Store, error) {
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
		log:  log,
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
		SELECT r.Id, r.State, r.Status
		FROM Rooms r 
		join RoomPlayers rp on rp.RoomID = r.Id 
		WHERE r.GameId = $1 and rp.PlayerId = $2 and not rp.IsQuit
		`, gameID, playerID)

	res := &model.Room{}
	err := row.Scan(&res.ID, &res.State, &res.Status)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return res, nil
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
		set Status = $2, State = $3
		where Id = $1
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

	err = s.storeRooms(ctx, rooms, stmtRoomInsert, stmtRoomUpdate, stmtRoomPlayerInsert)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *DBStore) storeRooms(
	ctx context.Context,
	rooms []*model.MatcherRoom,
	stmtRoomInsert *sql.Stmt,
	stmtRoomUpdate *sql.Stmt,
	stmtRoomPlayerInsert *sql.Stmt) error {
	for _, r := range rooms {
		if r.IsNew {
			var err error
			_, err = stmtRoomInsert.ExecContext(ctx, r.ID, r.GameID, r.Status, r.State)
			if err != nil {
				return err
			}
			s.log.Debug("room inserted", "ID", r.ID, "GameID", r.GameID, "Status", r.Status, "State", r.State)
		} else if r.StatusChanged {
			var err error
			_, err = stmtRoomUpdate.ExecContext(ctx, r.ID, r.Status, r.State)
			if err != nil {
				return err
			}
			s.log.Debug("room updated", "ID", r.ID, "GameID", r.GameID, "Status", r.Status, "State", r.State)
		}

		for _, p := range r.Players {
			var err error
			if p.IsNew {
				_, err = stmtRoomPlayerInsert.ExecContext(ctx, p.PlayerID, r.ID)
				if err != nil {
					return err
				}
				s.log.Debug("roomPlayer inserted", "PlayerID", p.PlayerID, "RoomID", r.ID)
			}
		}
	}

	return nil
}

func (s *DBStore) LoadWaitingRooms(ctx context.Context) ([]*model.MatcherRoom, error) {
	res := make([]*model.MatcherRoom, 0)
	var rows *sql.Rows

	// одним запросом загружаем комнаты и игроков
	rows, err := s.conn.QueryContext(ctx, `
		select r.Id, r.GameId, r.Status, rp.PlayerId
		from Rooms r
		left join RoomPlayers rp on rp.RoomId = r.Id
		where r.Status = 'wait' and not (rp.IsQuit is not null and rp.IsQuit)
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		curRoomID uuid.UUID
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
		if curRoomID != roomID {
			room = &model.MatcherRoom{
				ID:      roomID,
				Players: make([]*model.MatcherPlayer, 0),
				Status:  status,
				GameID:  gameID,
			}
			curRoomID = room.ID
			res = append(res, room)
		} else if playerID != nil {
			room.Players = append(room.Players, &model.MatcherPlayer{
				PlayerID: *playerID,
			})
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
