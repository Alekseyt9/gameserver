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

	var playerID string
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "playerID" {
			playerID = cookie.Value
			break
		}
	}
	require.True(t, playerID != "")

	req, err = http.NewRequest(http.MethodPost, ts.URL+"/api/room/connect", nil)
	require.NoError(t, err)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	resp, err = ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	require.True(t, bodyString != "")

}
