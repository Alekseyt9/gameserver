package test_test

import (
	"gameserver/internal/run"
	"gameserver/internal/services"
	"gameserver/internal/services/store"
	"gameserver/internal/test"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	ts    *httptest.Server
	store store.Store
}

func (suite *TestSuite) SetupSuite() {
	store := test.NewMemStore()
	suite.store = store
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	pm := services.NewPlayerManager(store)
	gm := services.NewGameManager(store, pm)
	m, err := services.NewMatcher(store, pm, gm, log)
	suite.Require().NoError(err)

	rm := services.NewRoomManager(store, gm, pm, m, log)

	cfg := &run.Config{}
	r := run.Router(store, pm, rm, cfg, log)
	server := httptest.NewServer(r)
	suite.ts = server
}

func (suite *TestSuite) TearDownSuite() {
	suite.ts.Close()
}

func TestRouterSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
