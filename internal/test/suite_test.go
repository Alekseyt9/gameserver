package test_test

import (
	"gameserver/internal/run"
	"gameserver/internal/services"
	"gameserver/internal/services/store"
	"gameserver/internal/test"
	"net/http/httptest"
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

	pm := services.NewPlayerManager(store)
	gm := services.NewGameManager(store, pm)
	m, err := services.NewMatcher(store, pm, gm)
	suite.Require().NoError(err)

	rm := services.NewRoomManager(store, gm, pm, m)

	cfg := &run.Config{}
	r := run.Router(store, pm, rm, cfg)
	server := httptest.NewServer(r)
	suite.ts = server
}

func (suite *TestSuite) TearDownSuite() {
	suite.ts.Close()
}

func TestRouterSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
