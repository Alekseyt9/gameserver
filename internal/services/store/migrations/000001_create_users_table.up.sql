
CREATE TABLE Players ( 
    PlayerId 	UUID PRIMARY KEY, 
    Name 		TEXT NOT NULL 
);

CREATE TABLE Rooms (
    RoomId 		    UUID PRIMARY KEY,
    GameType 		TEXT NOT NULL,
    State 		    TEXT NOT NULL,
    LastMove 		TIME NOT NULL,
    DeadlineMove 	TIME NOT NULL
);

CREATE TABLE RoomPlayers (
    PlayerId 	UUID,
    RoomId 	UUID,
    PRIMARY KEY (PlayerId, RoomId),
    FOREIGN KEY (PlayerId) REFERENCES Players(PlayerId),
    FOREIGN KEY (RoomId) REFERENCES Rooms(RoomId) ON DELETE CASCADE
);


