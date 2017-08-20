package model

type Permission struct {
	Id       uint64 `json:",omitempty"`
	User_id  uint64 `json:",omitempty"`
	Space_id uint64 `json:",omitempty"`
}
