
CREATE TABLE Players ( 
    Id 	        UUID PRIMARY KEY, 
    Name 		TEXT NOT NULL 
);

CREATE TABLE Rooms (
    Id 		        UUID PRIMARY KEY,
    GameId		    TEXT NOT NULL,
    State 		    TEXT NOT NULL,
    Status          TEXT NOT NULL
    --LastMove 		TIME NOT NULL,
    --DeadlineMove 	TIME NOT NULL
);
CREATE INDEX rooms_game_id ON Rooms (GameId);

CREATE TABLE RoomPlayers (
    PlayerId 	UUID,
    RoomId 	    UUID,
    IsQiut      boolean,    -- игрок вышел из комнаты
    PRIMARY KEY (PlayerId, RoomId),
    FOREIGN KEY (PlayerId) REFERENCES Players(Id),
    FOREIGN KEY (RoomId) REFERENCES Rooms(Id) ON DELETE CASCADE
);


