package game

type TTTMessage struct {
	Action string
	Data   string
}

type TTTMoveData struct {
	Move [2]byte
}
