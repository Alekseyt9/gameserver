package services

type RoomManager struct {
}

/*
Получает существующую комнату для игрока или создает новую
Существующую: если игра с таким типом для игрока уже есть или игры нет, но игрок подключается к свободной комнате
Новую: если нет свободной комнаты для игрока и нет существующей игры игры
*/
func (*RoomManager) GetOrCreateRoom(userId string, gameType string) {
	// запрос в хранилище по игроку и типу игры
	// ? нужна ли комната в базе, котоорая ожидает игроков?
	// наверное да, тк чтобы игроки могли подключаться, когда другие офлайн
}
