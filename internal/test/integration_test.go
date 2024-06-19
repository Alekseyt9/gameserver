package test_test

import (
	"bytes"
	json "encoding/json"
	game "gameserver/internal/services/game/tictactoe"
	"gameserver/internal/test"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	InitialState = iota
	RecieveGameDataState
	MakeMoveState
	QuitRoomState
)

func (suite *TestSuite) TestIntegration() {
	cookies1, playerID1 := playerRegister(suite)
	//ws1 := createWSDial(suite, cookies1)
	ws1 := connectToRoom(suite, cookies1)

	cookies2, playerID2 := playerRegister(suite)
	//ws2 := createWSDial(suite, cookies2)
	ws2 := connectToRoom(suite, cookies2)

	// процесс игры для 1го игрока.
	gameProcess(suite, ws1, playerID1, cookies1)

	// процесс игры для 2го игрока.
	gameProcess(suite, ws2, playerID2, cookies2)

	time.Sleep(time.Second * 10)
}

// процесс игры.
func gameProcess(suite *TestSuite, ws *websocket.Conn, playerID *uuid.UUID, cookies []*http.Cookie) {
	t := suite.T()

	go func() {
		state := InitialState

		for {
			require.NotNil(t, playerID)
			_, msg, err := ws.ReadMessage()
			require.NoError(t, err)
			require.NotEqual(t, "", string(msg))
			var m test.OutMessage
			err = m.UnmarshalJSON(msg)
			require.NoError(t, err)

			switch state {
			case InitialState:
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
					state = RecieveGameDataState
					continue
				}

			case RecieveGameDataState:
				var s game.TTTSendState
				var data []byte
				data, err = json.Marshal(m.Data.Data)
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
						}`))
					require.NoError(t, err)
					state = MakeMoveState
					continue
				}

			case MakeMoveState:
				var st game.TTTSendState
				var data []byte
				data, err = json.Marshal(m.Data.Data)
				require.NoError(t, err)
				err = st.UnmarshalJSON(data)
				require.NoError(t, err)
				quitRoom(suite, cookies)
				state = QuitRoomState
				err = ws.Close()
				require.NoError(t, err)
				continue
			}
		}
	}()
}

func quitRoom(suite *TestSuite, cookies []*http.Cookie) {
	ts := suite.ts
	t := suite.T()

	jsonValue := []byte(`{"gameID":"tictactoe"}`)
	req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/room/quit", bytes.NewBuffer(jsonValue))
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
	require.NotEqual(t, "", bodyString)
}

// подключение игрока к комнате.
func connectToRoom(suite *TestSuite, cookies []*http.Cookie) *websocket.Conn {
	//ts := suite.ts
	//t := suite.T()

	/*
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
		require.NotEqual(t, "", bodyString)
	*/

	return createWSDial(suite, cookies)
}

// создание websocket соединения.
func createWSDial(s *TestSuite, cookies []*http.Cookie) *websocket.Conn {
	ts := s.ts
	t := s.T()

	header := http.Header{}
	for _, cookie := range cookies {
		header.Add("Cookie", cookie.String())
	}
	wsURL := "ws" + ts.URL[len("http"):] + "/api/room/connect?gameid=tictactoe"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, header)
	assert.NoError(t, err)
	return ws
}

// регистрация пользователя.
func playerRegister(suite *TestSuite) ([]*http.Cookie, *uuid.UUID) {
	ts := suite.ts
	t := suite.T()

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
	require.NotEqual(t, "", playerID)
	pID, err := uuid.Parse(playerID)
	require.NoError(t, err)

	return cookies, &pID
}
