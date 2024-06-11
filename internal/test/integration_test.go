package test

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestIntegration() {
	//ts := s.ts
	t := s.T()

	cookies1 := playerRegister(s)
	ws1 := createWSDial(s, cookies1)
	connectToRoom(s, cookies1)

	cookies2 := playerRegister(s)
	ws2 := createWSDial(s, cookies2)
	connectToRoom(s, cookies2)

	// ожидание ответа о старте игры по websocket
	go func() {
		for {
			_, message, err := ws1.ReadMessage()
			require.NoError(t, err)
			require.True(t, string(message) != "")
		}
	}()

	go func() {
		for {
			_, message, err := ws2.ReadMessage()
			require.NoError(t, err)
			require.True(t, string(message) != "")
		}
	}()

	time.Sleep(time.Second * 60)

	// запрос состояния игры

	// делаем ход

	// тестим рассылку стейта игры

	// отключаемся от комнаты
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
