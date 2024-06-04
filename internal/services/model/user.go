package model

import "github.com/beevik/guid"

type Player struct {
	ID   guid.Guid
	Name string
}
