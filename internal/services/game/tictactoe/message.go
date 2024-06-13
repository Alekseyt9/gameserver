package game

//easyjson:json
type TTTMessage struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

//easyjson:json
type TTTMoveData struct {
	Move [2]int `json:"move"`
}
