package test

import (
	"io"
	"net/http"

	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestIntegration() {
	ts := s.ts
	t := s.T()

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/player/register", nil)
	require.NoError(t, err)
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	_, err = io.Copy(io.Discard, resp.Body)
	require.NoError(t, err)
}
