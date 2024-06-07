package store

import (
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

func (db *DBStore) GetUser(id guid.Guid) (*model.Player, error) {
	return nil, errors.New("not implemented")
}

func (db *DBStore) CreateUser(player *model.Player) error {
	return errors.New("not implemented")
}

func (db *DBStore) GetRoomState(playerID guid.Guid, gameType string) (string, error) {
	return "", errors.New("not implemented")
}

func (db *DBStore) SetRoomState(id guid.Guid, state string) error {
	return errors.New("not implemented")
}

func (db *DBStore) CreateRoom() error {
	return errors.New("not implemented")
}

func (db *DBStore) DropRoom(id guid.Guid) error {
	return errors.New("not implemented")
}

func (db *DBStore) CreateOrUpdateRooms([]model.MatcherRoom) error {
	return errors.New("not implemented")
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
