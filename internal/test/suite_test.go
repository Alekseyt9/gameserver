package test

import (
	"gameserver/internal/run"
	"gameserver/internal/services"
	"gameserver/internal/services/store"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	ts    *httptest.Server
	ws    *websocket.Conn
	store store.Store
}

func (suite *TestSuite) SetupSuite() {
	s := store.NewMemStore()
	suite.store = s

	pm := services.NewPlayerManager(s)
	gm := services.NewGameManager(s, pm)
	m, err := services.NewMatcher(s)
	assert.NoError(suite.T(), err)

	rm := services.NewRoomManager(s, gm, m)
	cfg := &run.Config{}

	r := run.Router(s, pm, rm, cfg)
	server := httptest.NewServer(r)
	suite.ts = server
}

func (suite *TestSuite) TearDownSuite() {
	suite.ts.Close()
	suite.ws.Close()
}

func TestRouterSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
