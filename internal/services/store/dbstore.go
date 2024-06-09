package store

import (
	"context"
	"database/sql"
	"errors"
	"gameserver/internal/services/model"
	"path/filepath"
	"runtime"

	"github.com/beevik/guid"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // needs for init
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

func (s *DBStore) GetUser(ctx context.Context, id guid.Guid) (*model.Player, error) {
	row := s.conn.QueryRowContext(ctx, `SELECT Name FROM Players WHERE Id = $1`, id)
	var name string
	err := row.Scan(&name)

	if err != nil {
		return nil, err
	}

	return &model.Player{ID: id, Name: name}, nil
}

func (s *DBStore) CreateUser(ctx context.Context, p *model.Player) error {
	_, err := s.conn.ExecContext(ctx, `
		insert into Players(Id, Name)
		values ($1, $2)
	`, p.ID, p.Name)
	return err
}

func (s *DBStore) GetRoom(ctx context.Context, playerID guid.Guid, gameType string) (*model.Room, error) {
	row := s.conn.QueryRowContext(ctx, `
	SELECT r.ID, r.State
	FROM Rooms r 
	join RoomsPlayers rp on rp.RoomID = r.ID 
	WHERE r.GameType = $1 and rp.ID = $2`, gameType, playerID)

	var id guid.Guid
	var state string
	err := row.Scan(&id, &state)
	if err != nil {
		return nil, err
	}

	return &model.Room{ID: id, State: state}, nil
}

func (s *DBStore) SetRoomState(ctx context.Context, id guid.Guid, state string) error {
	_, err := s.conn.ExecContext(ctx, `
	update Rooms
	set State = $1
	where Id = $2
	`, id, state)
	return err
}

func (s *DBStore) CreateRoom(ctx context.Context) error {
	return errors.New("not implemented")
}

func (s *DBStore) DropRoom(ctx context.Context, id guid.Guid) error {
	return errors.New("not implemented")
}

func (s *DBStore) CreateOrUpdateRooms(ctx context.Context, rooms []model.MatcherRoom) error {
	// TODO создавать ID для комнат и игроков

	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck //defer

	return tx.Commit()
}

func (s *DBStore) LoadWaitingRooms(ctx context.Context) ([]model.MatcherRoom, error) {
	res := make([]model.MatcherRoom, 0)
	var rows *sql.Rows

	rows, err := s.conn.QueryContext(ctx, "select Id, GameId, Status from Rooms where Status = 'wait'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// TODO одним запросом загружать игроков

	for rows.Next() {
		r := model.MatcherRoom{}
		if err = rows.Scan(&r.ID, &r.GameID, &r.Status); err != nil {
			return nil, err
		}
		res = append(res, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
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
