package model

import "github.com/google/uuid"

type Room struct {
	ID     uuid.UUID
	State  string
	Status string
}
