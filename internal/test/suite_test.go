package test

import (
	"gameserver/internal/run"
	"gameserver/internal/services"
	"gameserver/internal/services/store"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	ts    *httptest.Server
	store store.Store
}

func (suite *TestSuite) SetupSuite() {
	s := NewMemStore()
	suite.store = s

	pm := services.NewPlayerManager(s)
	gm := services.NewGameManager(s, pm)
	m, err := services.NewMatcher(s, pm, gm)
	assert.NoError(suite.T(), err)

	rm := services.NewRoomManager(s, gm, pm, m)
	cfg := &run.Config{}

	r := run.Router(s, pm, rm, cfg)
	server := httptest.NewServer(r)
	suite.ts = server
}

func (suite *TestSuite) TearDownSuite() {
	suite.ts.Close()
}

func TestRouterSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
