package test

import (
	"bytes"
	json "encoding/json"
	game "gameserver/internal/services/game/tictactoe"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestIntegration() {
	cookies1 := playerRegister(s)
	ws1 := createWSDial(s, cookies1)
	connectToRoom(s, cookies1)

	cookies2 := playerRegister(s)
	ws2 := createWSDial(s, cookies2)
	connectToRoom(s, cookies2)

	// процесс игры для 1го игрока
	gameProcess(s, ws1, 1)

	// процесс игры для 2го игрока
	gameProcess(s, ws2, 2)

	time.Sleep(time.Second * 60)
}

// процесс игры
func gameProcess(s *TestSuite, ws *websocket.Conn, pNum byte) {
	t := s.T()

	go func() {
		state := 0

		for {
			_, msg, err := ws.ReadMessage()
			require.NoError(t, err)
			require.True(t, string(msg) != "")
			var m OutMessage
			err = m.UnmarshalJSON(msg)
			require.NoError(t, err)

			switch state {
			case 0:
				if m.Data.Action == "start" {
					err = ws.WriteMessage(websocket.TextMessage, []byte(`
						{
							"type": "game",
							"gameid": "tictactoe",
							"data": { 				
								"action": "state"
							}
						}					
					`))
					require.NoError(t, err)
					state = 1
					continue
				}

			case 1:
				var s game.TTTSendState
				data, err := json.Marshal(m.Data.Data)
				require.NoError(t, err)
				err = s.UnmarshalJSON(data)
				require.NoError(t, err)
				if s.Turn == s.You {
					err = ws.WriteMessage(websocket.TextMessage, []byte(`
					{
						"type": "game",
						"gameid": "tictactoe",
						"data": { 	
							"action":"move",
							"data": {
								"move": [1, 1]
							}								
						}
					}					
				`))
					require.NoError(t, err)
					state = 2
					continue
				}

			case 2:
				var s game.TTTSendState
				data, err := json.Marshal(m.Data.Data)
				require.NoError(t, err)
				err = s.UnmarshalJSON(data)
				require.NoError(t, err)
			}
		}
	}()
}

// подключение игрока к комнате
func connectToRoom(s *TestSuite, cookies []*http.Cookie) {
	ts := s.ts
	t := s.T()

	jsonValue := []byte(`{"gameID":"tictactoe"}`)
	req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/room/connect", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, err)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	require.True(t, bodyString != "")
}

// создание websocket соединения
func createWSDial(s *TestSuite, cookies []*http.Cookie) *websocket.Conn {
	ts := s.ts
	t := s.T()

	header := http.Header{}
	for _, cookie := range cookies {
		header.Add("Cookie", cookie.String())
	}
	wsURL := "ws" + ts.URL[len("http"):] + "/ws"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, header)
	assert.NoError(t, err)
	return ws
}

// регистрация пользователя
func playerRegister(s *TestSuite) []*http.Cookie {
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

	return cookies
}
