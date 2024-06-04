package model

import "github.com/beevik/guid"

type User struct {
	ID   guid.Guid
	Name string
}
