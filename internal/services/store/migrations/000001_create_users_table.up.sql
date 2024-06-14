
CREATE TABLE Players ( 
    Id 	        UUID PRIMARY KEY, 
    Name 		TEXT NOT NULL 
);

CREATE TABLE Rooms (
    Id 		        UUID PRIMARY KEY,
    GameId		    TEXT NOT NULL,
    State 		    TEXT NOT NULL,
    Status          TEXT NOT NULL
);
CREATE INDEX i_rooms_game_id ON Rooms (GameId);

CREATE TABLE RoomPlayers (
    PlayerId 	UUID,
    RoomId 	    UUID,
    IsQuit      boolean DEFAULT false,    -- игрок вышел из комнаты
    FOREIGN KEY (PlayerId) REFERENCES Players(Id),
    FOREIGN KEY (RoomId) REFERENCES Rooms(Id) ON DELETE CASCADE
);
CREATE INDEX i_roomsplayers_playerid_roomid ON RoomPlayers(PlayerId, RoomId);
