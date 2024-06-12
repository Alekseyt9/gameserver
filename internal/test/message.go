package test

//easyjson:json
type OutMessage struct {
	Type   string        `json:"type"`
	GameID string        `json:"gameid"`
	Data   ActionMessage `json:"data"`
}

//easyjson:json
type ActionMessage struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
}
