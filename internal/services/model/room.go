package model

import "github.com/beevik/guid"

type Room struct {
	ID     guid.Guid
	State  string
	Status string
}
