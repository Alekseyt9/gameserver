package test

import (
	"gameserver/internal/run"
	"gameserver/internal/services"
	"gameserver/internal/services/store"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	ts *httptest.Server
}

func (suite *TestSuite) SetupSuite() {
	s := store.NewMemStore()
	pm := services.NewPlayerManager(s)
	gm := services.NewGameManager(s, pm)
	m, err := services.NewMatcher()
	if err != nil {
		// TODO log
	}
	rm := services.NewRoomManager(s, gm, m)
	cfg := &run.Config{}

	r := run.Router(s, pm, rm, cfg)
	suite.ts = httptest.NewServer(r)
}

func (suite *TestSuite) TearDownSuite() {
	suite.ts.Close()
}

func TestRouterSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
