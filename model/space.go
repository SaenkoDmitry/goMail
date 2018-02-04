package model

type Space struct {
	Id     uint64  `json:",omitempty"`
	Name  string `json:",omitempty"`
	UserId uint64  `json:",omitempty"`
}