package model

type User struct {
	Id           uint64  `json:",omitempty"`
	Name         string  `json:",omitempty"`
	HashPassword string  `json:"-"`
}
